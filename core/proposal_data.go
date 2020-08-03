package core

import (
	"bytes"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

// copy from prysm https://github.com/prysmaticlabs/prysm/blob/master/validator/client/validator_propose.go#L220-L226
type BeaconBlockHeader struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    []byte `ssz-size:"32"`
	StateRoot     []byte `ssz-size:"32"`
	BodyRoot      []byte `ssz-size:"32"`
}

func (proposal *BeaconBlockHeader) Compare(other *BeaconBlockHeader) bool {
	return proposal.Slot == other.Slot &&
		bytes.Compare(proposal.ParentRoot, other.ParentRoot) == 0 &&
		proposal.ProposerIndex == other.ProposerIndex &&
		bytes.Compare(proposal.StateRoot, other.StateRoot) == 0 &&
		bytes.Compare(proposal.BodyRoot, other.BodyRoot) == 0

}

func ToCoreBlockData(req *pb.SignBeaconProposalRequest) *BeaconBlockHeader {
	return &BeaconBlockHeader{ // Create a local copy of the data; we need ssz size information to calculate the correct root.
		Slot:          req.Data.Slot,
		ProposerIndex: req.Data.ProposerIndex,
		ParentRoot:    req.Data.ParentRoot,
		StateRoot:     req.Data.StateRoot,
		BodyRoot:      req.Data.BodyRoot,
	}
}

type ProposalDetectionType string

const (
	DoubleProposal ProposalDetectionType = "DoubleProposal"
	ValidProposal  ProposalDetectionType = "Valid"
	Error          ProposalDetectionType = "Error"
)

type ProposalSlashStatus struct {
	Proposal *BeaconBlockHeader
	Status   ProposalDetectionType
	Error    error
}
