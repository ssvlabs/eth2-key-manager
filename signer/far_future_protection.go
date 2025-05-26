package signer

import (
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// MaxFarFutureDelta is the max delta of far future signing
var MaxFarFutureDelta = time.Minute * 20

func IsValidFarFutureEpoch(network Network, epoch phase0.Epoch) bool {
	maxValidEpoch := network.EstimatedEpochAtSlot(network.EstimatedSlotAtTime(time.Now().Add(MaxFarFutureDelta)))
	return epoch <= maxValidEpoch
}

// IsValidFarFutureSlot returns true if the given slot is valid
func IsValidFarFutureSlot(network Network, slot phase0.Slot) bool {
	maxValidSlot := network.EstimatedSlotAtTime(time.Now().Add(MaxFarFutureDelta))
	return slot <= maxValidSlot
}
