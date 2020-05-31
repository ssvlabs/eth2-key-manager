package core

import (
	"bytes"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

// copy from prysm
type Checkpoint struct {
	Epoch uint64
	Root  []byte `ssz-size:"32"`
}

// returns true if equal, false if not
func (checkpoint *Checkpoint) compare (other *Checkpoint) bool {
	return checkpoint.Epoch == other.Epoch && bytes.Compare(checkpoint.Root, other.Root) == 0
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

// returns true if equal, false if not
func (att *BeaconAttestation) compare(other *BeaconAttestation) bool {
	return att.Slot == other.Slot &&
		att.CommitteeIndex == other.CommitteeIndex &&
		bytes.Compare(att.BeaconBlockRoot, other.BeaconBlockRoot) == 0 &&
		att.Target.compare(other.Target) &&
		att.Source.compare(other.Source)

}

// will return an array of attestations that this attestation will slash based on a provided history
// will detect double, surround and surrounded slashable events
func (att *BeaconAttestation) SlashesAttestations (history []*BeaconAttestation) []*BeaconAttestation {
	ret := make ([]*BeaconAttestation,0)

	if val := detectDoubleVote(att,history); val != nil {
		ret = append(ret,val)
	}

	return ret
}

func detectDoubleVote(att *BeaconAttestation, history []*BeaconAttestation) *BeaconAttestation {
	for _,history_att := range history {
		if att.Target.Epoch == history_att.Target.Epoch {
			if !att.compare(history_att) {
				return history_att
			}
		}
	}
	return nil
}