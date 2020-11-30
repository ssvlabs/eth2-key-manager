package core

import (
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

type ProposalDetectionType string

const (
	DoubleProposal ProposalDetectionType = "DoubleProposal"
	ValidProposal  ProposalDetectionType = "Valid"
	Error          ProposalDetectionType = "Error"
)

type ProposalSlashStatus struct {
	Proposal *eth.BeaconBlock
	Status   ProposalDetectionType
	Error    error
}
