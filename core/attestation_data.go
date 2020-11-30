package core

import (
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

type VoteDetectionType string

const (
	DoubleVote             VoteDetectionType = "DoubleVote"
	SurroundingVote        VoteDetectionType = "SurroundingVote"
	SurroundedVote         VoteDetectionType = "SurroundedVote"
	HighestAttestationVote VoteDetectionType = "HighestAttestationVote"
)

type AttestationSlashStatus struct {
	Attestation *eth.AttestationData
	Status      VoteDetectionType
}
