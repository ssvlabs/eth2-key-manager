package core

import (
	"bytes"

	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

// copy from prysm
type Checkpoint struct {
	Epoch uint64 `json:"epoch"`
	Root  []byte `ssz-size:"32" json:"root"`
}

// returns true if equal, false if not
func (checkpoint *Checkpoint) compare(other *Checkpoint) bool {
	return checkpoint.Epoch == other.Epoch && bytes.Compare(checkpoint.Root, other.Root) == 0
}

// copy from prysm
type BeaconAttestation struct {
	Slot            uint64      `json:"slot"`
	CommitteeIndex  uint64      `json:"committee_index"`
	BeaconBlockRoot []byte      `ssz-size:"32" json:"beacon_block_root"`
	Source          *Checkpoint `json:"source"`
	Target          *Checkpoint `json:"target"`
}

type VoteDetectionType string

const (
	DoubleVote      VoteDetectionType = "DoubleVote"
	SurroundingVote VoteDetectionType = "SurroundingVote"
	SurroundedVote  VoteDetectionType = "SurroundedVote"
)

type AttestationSlashStatus struct {
	Attestation *BeaconAttestation
	Status      VoteDetectionType
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
func (att *BeaconAttestation) Compare(other *BeaconAttestation) bool {
	return att.Slot == other.Slot &&
		att.CommitteeIndex == other.CommitteeIndex &&
		bytes.Compare(att.BeaconBlockRoot, other.BeaconBlockRoot) == 0 &&
		att.Target.compare(other.Target) &&
		att.Source.compare(other.Source)

}

// will return an array of attestations that this attestation will slash based on a provided history
// will detect double, surround and surrounded slashable events
func (att *BeaconAttestation) SlashesAttestations(history []*BeaconAttestation) []*AttestationSlashStatus {
	ret := make([]*AttestationSlashStatus, 0)

	for _, history_att := range history {
		if val := detectDoubleVote(att, history_att); val != nil {
			ret = append(ret, &AttestationSlashStatus{
				Attestation: val,
				Status:      DoubleVote,
			})
		}

		// Surrounding vote
		if val := detectSurroundingVote(att, history_att); val != nil {
			ret = append(ret, &AttestationSlashStatus{
				Attestation: val,
				Status:      SurroundingVote,
			})
		}

		// Surrounded vote
		if val := detectSurroundedVote(att, history_att); val != nil {
			ret = append(ret, &AttestationSlashStatus{
				Attestation: val,
				Status:      SurroundedVote,
			})
		}
	}

	return ret
}

// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#is_slashable_attestation_data
func detectDoubleVote(att *BeaconAttestation, other *BeaconAttestation) *BeaconAttestation {
	if att.Target.Epoch == other.Target.Epoch {
		if !att.Compare(other) {
			return other
		}
	}
	return nil
}

// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#is_slashable_attestation_data
func detectSurroundingVote(att *BeaconAttestation, other *BeaconAttestation) *BeaconAttestation {
	if att.Source.Epoch < other.Source.Epoch && att.Target.Epoch > other.Target.Epoch {
		return other
	}
	return nil
}

func detectSurroundedVote(att *BeaconAttestation, other *BeaconAttestation) *BeaconAttestation {
	return detectSurroundingVote(other, att)
}
