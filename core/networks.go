package core

import (
	"encoding/hex"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Network represents the network.
type Network string

// NetworkFromString returns network from the given string value
func NetworkFromString(n string) (Network, error) {
	switch n {
	case string(PyrmontNetwork):
		return PyrmontNetwork, nil
	case string(PraterNetwork):
		return PraterNetwork, nil
	case string(MainNetwork):
		return MainNetwork, nil
	default:
		return "", errors.New("undefined network")
	}
}

// GenesisForkVersion returns the genesis fork version of the network.
func (n Network) GenesisForkVersion() phase0.Version {
	switch n {
	case PyrmontNetwork:
		return phase0.Version{0, 0, 32, 9}
	case PraterNetwork:
		return phase0.Version{0x00, 0x00, 0x10, 0x20}
	case MainNetwork:
		return phase0.Version{0, 0, 0, 0}
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return phase0.Version{}
	}
}

// GenesisValidatorsRoot returns the genesis validators root of the network.
func (n Network) GenesisValidatorsRoot() (phase0.Root, error) {
	var genValidatorsRoot phase0.Root
	switch n {
	case PraterNetwork:
		rootBytes, err := hex.DecodeString("043db0d9a83813551ee2f33450d23797757d430911a9320530ad8a0eabc43efb")
		if err != nil {
			return phase0.Root{}, err
		}
		copy(genValidatorsRoot[:], rootBytes)
	case MainNetwork:
		rootBytes, err := hex.DecodeString("4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95")
		if err != nil {
			return phase0.Root{}, err
		}
		copy(genValidatorsRoot[:], rootBytes)
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return phase0.Root{}, errors.New("undefined network")
	}
	return genValidatorsRoot, nil
}

// DepositContractAddress returns the deposit contract address of the network.
func (n Network) DepositContractAddress() string {
	switch n {
	case PyrmontNetwork:
		return "0x8c5fecdC472E27Bc447696F431E425D02dd46a8c"
	case PraterNetwork:
		return "0xff50ed3d0ec03ac01d4c79aad74928bff48a7b2b"
	case MainNetwork:
		return "0x00000000219ab540356cBB839Cbe05303d7705Fa"
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return ""
	}
}

// FullPath returns the full path of the network.
func (n Network) FullPath(relativePath string) string {
	return BaseEIP2334Path + relativePath
}

// MinGenesisTime returns min genesis time value
func (n Network) MinGenesisTime() uint64 {
	switch n {
	case PyrmontNetwork:
		return 1605700807
	case PraterNetwork:
		return 1616508000
	case MainNetwork:
		return 1606824023
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return 0
	}
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

// Available networks.
const (
	// PyrmontNetwork represents the Pyrmont test network.
	PyrmontNetwork Network = "pyrmont"

	// PraterNetwork represents the Prater test network.
	PraterNetwork Network = "prater"

	// MainNetwork represents the main network.
	MainNetwork Network = "mainnet"
)
