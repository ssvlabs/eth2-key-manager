package core

import (
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

type ProposalDetectionType string

const (
	DoubleProposal      ProposalDetectionType = "DoubleProposal"
	HighestProposalVote ProposalDetectionType = "HighestProposalVote"
	ValidProposal       ProposalDetectionType = "Valid"
	Error               ProposalDetectionType = "Error"
)

type ProposalSlashStatus struct {
	Proposal *eth.BeaconBlock
	Status   ProposalDetectionType
	Error    error
}
