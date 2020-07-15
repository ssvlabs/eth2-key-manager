package slashing_protection

import (
	"github.com/bloxapp/KeyVault/core"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
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
func (protector *NormalProtection) IsSlashableAttestation(account core.Account, req *pb.SignBeaconAttestationRequest) ([]*core.AttestationSlashStatus, error) {
	data := core.ToCoreAttestationData(req)

	lookupStartEpoch := lookupEpochSub(data.Source.Epoch, epochLookback)
	lookupEndEpoch := req.Data.Target.Epoch

	// lookupEndEpoch should be the latest written attestation, if not than req.Data.Target.Epoch
	latestAtt, err := protector.RetrieveLatestAttestation(account)
	if err != nil {
		return nil, err
	}
	if latestAtt != nil {
		lookupEndEpoch = latestAtt.Target.Epoch
	}

	history, err := protector.store.ListAttestations(account, lookupStartEpoch, lookupEndEpoch)
	if err != nil {
		return nil, err
	}

	return data.SlashesAttestations(history), nil
}

func (protector *NormalProtection) IsSlashableProposal(account core.Account, req *pb.SignBeaconProposalRequest) (*core.ProposalSlashStatus, error) {
	matchedProposal, err := protector.store.RetrieveProposal(account, req.Data.Slot)
	if err != nil && err.Error() != "proposal not found" {
		return nil, err
	}

	if matchedProposal == nil {
		return nil, nil
	}

	data := core.ToCoreBlockData(req)

	// if it's the same
	if data.Compare(matchedProposal) {
		return nil, nil
	}

	// slashable
	return &core.ProposalSlashStatus{
		Proposal: data,
		Status:   core.DoubleProposal,
	}, nil
}

func (protector *NormalProtection) SaveAttestation(account core.Account, req *pb.SignBeaconAttestationRequest) error {
	data := core.ToCoreAttestationData(req)
	err := protector.store.SaveAttestation(account, data)
	if err != nil {
		return err
	}
	return protector.SaveLatestAttestation(account, req)
}

func (protector *NormalProtection) SaveProposal(account core.Account, req *pb.SignBeaconProposalRequest) error {
	data := core.ToCoreBlockData(req)
	return protector.store.SaveProposal(account, data)
}

func (protector *NormalProtection) SaveLatestAttestation(account core.Account, req *pb.SignBeaconAttestationRequest) error {
	val, err := protector.store.RetrieveLatestAttestation(account)
	if err != nil {
		return nil
	}

	data := core.ToCoreAttestationData(req)
	if val == nil {
		return protector.store.SaveLatestAttestation(account, data)
	}
	if val.Target.Epoch < req.Data.Target.Epoch { // only write newer
		return protector.store.SaveLatestAttestation(account, data)
	}

	return nil
}

func (protector *NormalProtection) RetrieveLatestAttestation(account core.Account) (*core.BeaconAttestation, error) {
	return protector.store.RetrieveLatestAttestation(account)
}

// specialized func that will prevent overflow for lookup epochs for uint64
func lookupEpochSub(l uint64, r uint64) uint64 {
	if l >= r {
		return l - r
	}
	return 0
}
