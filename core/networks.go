package core

import (
	"time"

	"github.com/prysmaticlabs/prysm/shared/timeutils"
	"github.com/sirupsen/logrus"
)

// Network represents the network.
type Network string

// NetworkFromString returns network from the given string value
func NetworkFromString(n string) Network {
	switch n {
	case string(PyrmontNetwork):
		return PyrmontNetwork
	case string(MainNetwork):
		return MainNetwork
	default:
		return ""
	}
}

// ForkVersion returns the fork version of the network.
func (n Network) ForkVersion() []byte {
	switch n {
	case PyrmontNetwork:
		return []byte{0, 0, 32, 9}
	case MainNetwork:
		return []byte{0, 0, 0, 0}
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return nil
	}
}

// DepositContractAddress returns the deposit contract address of the network.
func (n Network) DepositContractAddress() string {
	switch n {
	case PyrmontNetwork:
		return "0x8c5fecdC472E27Bc447696F431E425D02dd46a8c"
	case MainNetwork:
		return "0x00000000219ab540356cBB839Cbe05303d7705Fa"
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return ""
	}
}

// ForkVersion returns the fork version of the network.
func (n Network) FullPath(relativePath string) string {
	return BaseEIP2334Path + relativePath
}

func (n Network) MinGenesisTime() uint64 {
	switch n {
	case PyrmontNetwork:
		return 1605700807
	case MainNetwork:
		return 1606824023
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return 0
	}
}

func (n Network) SlotDurationSec() time.Duration {
	return 12 * time.Second
}

func (n Network) SlotsPerEpoch() uint64 {
	return 32
}

func (n Network) EstimatedCurrentSlot() uint64 {
	return n.EstimatedSlotAtTime(timeutils.Now().Unix())
}

func (n Network) EstimatedSlotAtTime(time int64) uint64 {
	genesis := int64(n.MinGenesisTime())
	if time < genesis {
		return 0
	}
	return uint64(time-genesis) / uint64(n.SlotDurationSec().Seconds())
}

// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#compute_start_slot_at_epoch
func (n Network) EstimatedCurrentEpoch() uint64 {
	return n.EstimatedEpochAtSlot(n.EstimatedCurrentSlot())
}

func (n Network) EstimatedEpochAtSlot(slot uint64) uint64 {
	return slot / n.SlotsPerEpoch()
}

// Available networks.
const (
	// PyrmontNetwork represents the Pyrmont test network.
	PyrmontNetwork Network = "pyrmont"

	// MainNetwork represents the main network.
	MainNetwork Network = "mainnet"
)
