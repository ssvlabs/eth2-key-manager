package core

import (
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// ProposalDetectionType represents proposal slashing detection type
type ProposalDetectionType string

// Proposal slashing detection types
const (
	DoubleProposal      ProposalDetectionType = "DoubleProposal"
	HighestProposalVote ProposalDetectionType = "HighestProposalVote"
	ValidProposal       ProposalDetectionType = "Valid"
	Error               ProposalDetectionType = "Error"
)

// ProposalSlashStatus represents proposal slashing status
type ProposalSlashStatus struct {
	Proposal *eth.BeaconBlock
	Status   ProposalDetectionType
	Error    error
}
