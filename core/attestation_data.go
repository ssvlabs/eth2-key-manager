package core

import pb "github.com/wealdtech/eth2-signer-api/pb/v1"

// copy from prysm
type Checkpoint struct {
	Epoch uint64
	Root  []byte `ssz-size:"32"`
}

// copy from prysm
type BeaconAttestation struct {
	Slot            uint64
	CommitteeIndex  uint64
	BeaconBlockRoot []byte `ssz-size:"32"`
	Source          *Checkpoint
	Target          *Checkpoint
}

func ToCoreAttestationData(req *pb.SignBeaconAttestationRequest) *BeaconAttestation {
	return &BeaconAttestation{ // Create a local copy of the data; we need ssz size information to calculate the correct root.
		Slot:            req.Data.Slot,
		CommitteeIndex:  req.Data.CommitteeIndex,
		BeaconBlockRoot: req.Data.BeaconBlockRoot,
		Source: &Checkpoint{
			Epoch: req.Data.Source.Epoch,
			Root:  req.Data.Source.Root,
		},
		Target: &Checkpoint{
			Epoch: req.Data.Target.Epoch,
			Root:  req.Data.Target.Root,
		},
	}
}