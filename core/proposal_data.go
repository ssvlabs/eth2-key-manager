package core

import pb "github.com/wealdtech/eth2-signer-api/pb/v1"

// copy from prysm https://github.com/prysmaticlabs/prysm/blob/master/validator/client/validator_propose.go#L220-L226
type BeaconBlockHeader struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    []byte `ssz-size:"32"`
	StateRoot     []byte `ssz-size:"32"`
	BodyRoot      []byte `ssz-size:"32"`
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