package slashing_protection

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

type NormalProtection struct {
	store core.SlashingStore
}

func NewNormalProtection(store core.SlashingStore) *NormalProtection {
	return &NormalProtection{store: store}
}

// From prysm:
// We look back 128 epochs when updating min/max spans
// for incoming attestations.
// TODO - verify this is true
const epochLookback = 128

// will detect double, surround and surrounded slashable events
func (protector *NormalProtection) IsSlashableAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) ([]*core.AttestationSlashStatus, error) {
	data := core.ToCoreAttestationData(req)

	lookupStartEpoch := lookupEpochSub(data.Source.Epoch, epochLookback)
	lookupEndEpoch := req.Data.Target.Epoch

	// lookupEndEpoch should be the latest written attestation, if not than req.Data.Target.Epoch
	latestAtt, err := protector.RetrieveLatestAttestation(key)
	if err != nil {
		return nil, err
	}
	if latestAtt != nil {
		lookupEndEpoch = latestAtt.Target.Epoch
	}

	history, err := protector.store.ListAttestations(key, lookupStartEpoch, lookupEndEpoch)
	if err != nil {
		return nil, err
	}

	return data.SlashesAttestations(history), nil
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

func (protector *NormalProtection) SaveAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error {
	data := core.ToCoreAttestationData(req)
	err := protector.store.SaveAttestation(key, data)
	if err != nil {
		return err
	}
	return protector.SaveLatestAttestation(key, req)
}

func (protector *NormalProtection) SaveProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) error {
	data := core.ToCoreBlockData(req)
	return protector.store.SaveProposal(key, data)
}

func (protector *NormalProtection) SaveLatestAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error {
	val, err := protector.store.RetrieveLatestAttestation(key)
	if err != nil {
		return nil
	}

	data := core.ToCoreAttestationData(req)
	if val == nil {
		return protector.store.SaveLatestAttestation(key, data)
	}
	if val.Target.Epoch < req.Data.Target.Epoch { // only write newer
		return protector.store.SaveLatestAttestation(key, data)
	}

	return nil
}

func (protector *NormalProtection) RetrieveLatestAttestation(key e2types.PublicKey) (*core.BeaconAttestation, error) {
	return protector.store.RetrieveLatestAttestation(key)
}

// specialized func that will prevent overflow for lookup epochs for uint64
func lookupEpochSub(l uint64, r uint64) uint64 {
	if l >= r {
		return l - r
	}
	return 0
}
