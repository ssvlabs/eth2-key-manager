package slashing_protection

import (
	"fmt"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

type NormalProtection struct {
	store core.SlashingStore
}

// NewNormalProtection is the constructor of NormalProtection
func NewNormalProtection(store core.SlashingStore) *NormalProtection {
	return &NormalProtection{store: store}
}

// will detect double, surround and surrounded slashable events
func (protector *NormalProtection) IsSlashableAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) (*core.AttestationSlashStatus, error) {
	data := core.ToCoreAttestationData(req)

	// lookupEndEpoch should be the latest written attestation, if not than req.Data.Target.Epoch
	highest, err := protector.RetrieveHighestAttestation(key)
	if err != nil {
		return nil, err
	}
	if highest != nil {
		return highest.SlashesHighestAttestation(data), nil
	} else {
		return nil, fmt.Errorf("highest attestationn data is nil, can't determine if attestation is slashable")
	}
}

func (protector *NormalProtection) IsSlashableProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) *core.ProposalSlashStatus {
	matchedProposal, err := protector.store.RetrieveProposal(key, req.Data.Slot)
	if err != nil && err.Error() != "proposal not found" {
		return &core.ProposalSlashStatus{
			Proposal: nil,
			Status:   core.Error,
			Error:    err,
		}
	}

	if matchedProposal == nil {
		return &core.ProposalSlashStatus{
			Proposal: nil,
			Status:   core.ValidProposal,
		}
	}

	data := core.ToCoreBlockData(req)

	// if it's the same
	if data.Compare(matchedProposal) {
		return &core.ProposalSlashStatus{
			Proposal: data,
			Status:   core.ValidProposal,
		}
	}

	// slashable
	return &core.ProposalSlashStatus{
		Proposal: data,
		Status:   core.DoubleProposal,
	}
}

// Will potentially update the highest attestation given this latest attestation.
func (protector *NormalProtection) UpdateLatestAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error {
	data := core.ToCoreAttestationData(req)

	// if no previous highest attestation found, set current
	highest := protector.store.RetrieveHighestAttestation(key)
	if highest == nil {
		err := protector.store.SaveHighestAttestation(key, data)
		if err != nil {
			return err
		}
		return nil
	}

	// Taken from https://github.com/prysmaticlabs/prysm/blob/master/slasher/detection/detect.go#L233
	shouldUpdate := false
	if highest.Source.Epoch < data.Source.Epoch {
		highest.Source.Epoch = data.Source.Epoch
		shouldUpdate = true
	}
	if highest.Target.Epoch < data.Target.Epoch {
		highest.Target.Epoch = data.Target.Epoch
		shouldUpdate = true
	}

	if shouldUpdate {
		err := protector.store.SaveHighestAttestation(key, highest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (protector *NormalProtection) SaveProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) error {
	data := core.ToCoreBlockData(req)
	return protector.store.SaveProposal(key, data)
}

//func (protector *NormalProtection) SaveLatestAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error {
//	val, err := protector.store.RetrieveLatestAttestation(key)
//	if err != nil {
//		return nil
//	}
//
//	data := core.ToCoreAttestationData(req)
//	if val == nil {
//		return protector.store.SaveLatestAttestation(key, data)
//	}
//	if val.Target.Epoch < req.Data.Target.Epoch { // only write newer
//		return protector.store.SaveLatestAttestation(key, data)
//	}
//
//	return nil
//}

func (protector *NormalProtection) RetrieveHighestAttestation(key e2types.PublicKey) (*core.BeaconAttestation, error) {
	return protector.store.RetrieveHighestAttestation(key), nil
}

// specialized func that will prevent overflow for lookup epochs for uint64
func lookupEpochSub(l uint64, r uint64) uint64 {
	if l >= r {
		return l - r
	}
	return 0
}
