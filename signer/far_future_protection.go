package signer

import (
	"time"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/prysmaticlabs/prysm/shared/timeutils"
)

var FarFutureMaxValidEpoch = int64(time.Hour.Hours()) * 6

// To prevent far into the future signing request, verify a slot is within the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/validator.md#protection-best-practices
func IsValidFarFutureEpoch(network core.Network, epoch uint64) bool {
	maxValidEpoch := network.EstimatedEpochAtSlot(network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch))
	return epoch <= maxValidEpoch
}

func IsValidFarFutureSlot(network core.Network, slot uint64) bool {
	maxValidSlot := network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch)
	return slot <= maxValidSlot
}
