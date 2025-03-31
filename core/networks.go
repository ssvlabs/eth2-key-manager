package core

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/sirupsen/logrus"
)

// Network represents the network.
type Network string

// NetworkConfig stores configuration specific to an Ethereum network.
type NetworkConfig struct {
	GenesisForkVersion     phase0.Version
	GenesisValidatorsRoot  string
	DepositContractAddress string
	MinGenesisTime         uint64
}

// Available networks.
const (
	// PraterNetwork represents the Prater test network.
	PraterNetwork Network = "prater"

	// SepoliaNetwork represents the Sepolia test network.
	SepoliaNetwork Network = "sepolia"

	// HoleskyNetwork represents the Holesky test network.
	HoleskyNetwork Network = "holesky"

	// HoodiNetwork represents the Hoodi test network.
	HoodiNetwork Network = "hoodi"

	// MainNetwork represents the main network.
	MainNetwork Network = "mainnet"
)

// Network configurations.
var networks = map[Network]NetworkConfig{
	PraterNetwork: {
		GenesisForkVersion:     phase0.Version{0x00, 0x00, 0x10, 0x20},
		GenesisValidatorsRoot:  "043db0d9a83813551ee2f33450d23797757d430911a9320530ad8a0eabc43efb",
		DepositContractAddress: "0xff50ed3d0ec03ac01d4c79aad74928bff48a7b2b",
		MinGenesisTime:         1616508000,
	},
	SepoliaNetwork: {
		GenesisForkVersion:     phase0.Version{0x90, 0x00, 0x00, 0x69},
		GenesisValidatorsRoot:  "d8ea171f3c94aea21ebc42a1ed61052acf3f9209c00e4efbaaddac09ed9b8078",
		DepositContractAddress: "0x4242424242424242424242424242424242424242",
		MinGenesisTime:         1655733600,
	},
	HoleskyNetwork: {
		GenesisForkVersion:     phase0.Version{0x01, 0x01, 0x70, 0x00},
		GenesisValidatorsRoot:  "9143aa7c615a7f7115e2b6aac319c03529df8242ae705fba9df39b79c59fa8b1",
		DepositContractAddress: "0x4242424242424242424242424242424242424242",
		MinGenesisTime:         1695902400,
	},
	HoodiNetwork: {
		GenesisForkVersion:     phase0.Version{0x10, 0x00, 0x09, 0x10},
		GenesisValidatorsRoot:  "212f13fc4df078b6cb7db228f1c8307566dcecf900867401a92023d7ba99cb5f",
		DepositContractAddress: "0x00000000219ab540356cBB839Cbe05303d7705Fa",
		MinGenesisTime:         1742213400,
	},
	MainNetwork: {
		GenesisForkVersion:     phase0.Version{0, 0, 0, 0},
		GenesisValidatorsRoot:  "4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95",
		DepositContractAddress: "0x00000000219ab540356cBB839Cbe05303d7705Fa",
		MinGenesisTime:         1606824023,
	},
}

// NetworkFromString converts a string to a Network type.
func NetworkFromString(n string) Network {
	_, ok := networks[Network(n)]
	if !ok {
		return ""
	}

	return Network(n)
}

// NetworkFromForkVersion returns network from the given fork version
func NetworkFromForkVersion(version phase0.Version) (Network, error) {
	for net, cfg := range networks {
		if cfg.GenesisForkVersion == version {
			return net, nil
		}
	}
	return "", fmt.Errorf("network not found for the given fork version")
}

// GenesisForkVersion returns the genesis fork version of the network.
func (n Network) GenesisForkVersion() phase0.Version {
	if cfg, exists := networks[n]; exists {
		return cfg.GenesisForkVersion
	}
	logrus.WithField("network", n).Fatal("undefined network")
	return phase0.Version{}
}

// GenesisValidatorsRoot returns the genesis validators root of the network.
func (n Network) GenesisValidatorsRoot() phase0.Root {
	var root phase0.Root
	if cfg, exists := networks[n]; exists {
		if cfg.GenesisValidatorsRoot == "" {
			return root
		}
		rootBytes, err := hex.DecodeString(cfg.GenesisValidatorsRoot)
		if err != nil {
			logrus.WithError(err).Fatal("invalid genesis validators root")
		}
		copy(root[:], rootBytes)
		return root
	}
	logrus.WithField("network", n).Fatal("undefined network")
	return root
}

// DepositContractAddress returns the deposit contract address of the network.
func (n Network) DepositContractAddress() string {
	if cfg, exists := networks[n]; exists {
		return cfg.DepositContractAddress
	}
	logrus.WithField("network", n).Fatal("undefined network")
	return ""
}

// MinGenesisTime returns the min genesis time of the network.
func (n Network) MinGenesisTime() uint64 {
	if cfg, exists := networks[n]; exists {
		return cfg.MinGenesisTime
	}
	logrus.WithField("network", n).Fatal("undefined network")
	return 0
}

// FullPath returns the full path of the network.
func (n Network) FullPath(relativePath string) string {
	return BaseEIP2334Path + relativePath
}

// SlotDurationSec returns slot duration
func (n Network) SlotDurationSec() time.Duration {
	return 12 * time.Second
}

// SlotsPerEpoch returns number of slots per one epoch
func (n Network) SlotsPerEpoch() uint64 {
	return 32
}

// EstimatedCurrentSlot returns the estimation of the current slot
func (n Network) EstimatedCurrentSlot() phase0.Slot {
	return n.EstimatedSlotAtTime(time.Now().Unix())
}

// EstimatedSlotAtTime estimates slot at the given time
func (n Network) EstimatedSlotAtTime(time int64) phase0.Slot {
	genesis := int64(n.MinGenesisTime())
	if time < genesis {
		return 0
	}
	return phase0.Slot(uint64(time-genesis) / uint64(n.SlotDurationSec().Seconds()))
}

// EstimatedCurrentEpoch estimates the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#compute_start_slot_at_epoch
func (n Network) EstimatedCurrentEpoch() phase0.Epoch {
	return n.EstimatedEpochAtSlot(n.EstimatedCurrentSlot())
}

// EstimatedEpochAtSlot estimates epoch at the given slot
func (n Network) EstimatedEpochAtSlot(slot phase0.Slot) phase0.Epoch {
	return phase0.Epoch(slot / phase0.Slot(n.SlotsPerEpoch()))
}
