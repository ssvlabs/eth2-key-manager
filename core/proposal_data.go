package core

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
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
	Slot   phase0.Slot
	Status ProposalDetectionType
	Error  error
}
