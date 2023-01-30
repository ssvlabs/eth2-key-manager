package core

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VoteDetectionType represents vote detection type
type VoteDetectionType string

// Vote detection types
const (
	DoubleVote             VoteDetectionType = "DoubleVote"
	SurroundingVote        VoteDetectionType = "SurroundingVote"
	SurroundedVote         VoteDetectionType = "SurroundedVote"
	HighestAttestationVote VoteDetectionType = "HighestAttestationVote"
)

// AttestationSlashStatus represents attestation slashing status
type AttestationSlashStatus struct {
	Attestation *phase0.AttestationData
	Status      VoteDetectionType
}
