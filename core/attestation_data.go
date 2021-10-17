package core

import (
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
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
	Attestation *eth.AttestationData
	Status      VoteDetectionType
}
