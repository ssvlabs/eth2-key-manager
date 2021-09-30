package signer

import (
	"time"

	types "github.com/prysmaticlabs/eth2-types"

	"github.com/prysmaticlabs/prysm/shared/timeutils"

	"github.com/bloxapp/eth2-key-manager/core"
)

// FarFutureMaxValidEpoch is the max epoch of far future signing
var FarFutureMaxValidEpoch = int64(time.Minute.Seconds() * 20)

// IsValidFarFutureEpoch prevents far into the future signing request, verify a slot is within the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/validator.md#protection-best-practices
func IsValidFarFutureEpoch(network core.Network, epoch types.Epoch) bool {
	maxValidEpoch := network.EstimatedEpochAtSlot(network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch))
	return epoch <= maxValidEpoch
}

// IsValidFarFutureSlot returns true if the given slot is valid
func IsValidFarFutureSlot(network core.Network, slot types.Slot) bool {
	maxValidSlot := network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch)
	return slot <= maxValidSlot
}
