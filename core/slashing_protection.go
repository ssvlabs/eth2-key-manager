package core

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

type VaultSlashingProtector interface {
	IsSlashableAttestation(account Account, req *pb.SignBeaconAttestationRequest) ([]*AttestationSlashStatus,error)
	IsSlashableProposal(account Account, req *pb.SignBeaconProposalRequest) (*ProposalSlashStatus,error)
	SaveAttestation(account Account, req *pb.SignBeaconAttestationRequest) error
	SaveProposal(account Account, req *pb.SignBeaconProposalRequest) error
	SaveLatestAttestation(account Account, req *pb.SignBeaconAttestationRequest) error
	RetrieveLatestAttestation(account Account) (*BeaconAttestation, error)
}

type SlashingStore interface {
	SaveAttestation(account Account, req *BeaconAttestation) error
	RetrieveAttestation(account Account, epoch uint64) (*BeaconAttestation, error)
	// both epochStart and epochEnd reflect saved attestations by their target epoch
	ListAttestations(account Account, epochStart uint64, epochEnd uint64) ([]*BeaconAttestation, error)
	SaveProposal(account Account, req *BeaconBlockHeader) error
	RetrieveProposal(account Account, slot uint64) (*BeaconBlockHeader, error)
	SaveLatestAttestation(account Account, req *BeaconAttestation) error
	RetrieveLatestAttestation(account Account) (*BeaconAttestation, error)
}

type NormalProtection struct {
	store SlashingStore
}

func NewNormalProtection(store SlashingStore) *NormalProtection {
	return &NormalProtection{store:store}
}

// From prysm:
// We look back 128 epochs when updating min/max spans
// for incoming attestations.
// TODO - verify this is true
const epochLookback = 128

// will detect double, surround and surrounded slashable events
func (protector *NormalProtection) IsSlashableAttestation(account Account, req *pb.SignBeaconAttestationRequest) ([]*AttestationSlashStatus,error) {
	data := ToCoreAttestationData(req)

	lookupStartEpoch := lookupEpochSub(data.Source.Epoch, epochLookback)
	lookupEndEpoch := req.Data.Target.Epoch

	// lookupEndEpoch should be the latest written attestation, if not than req.Data.Target.Epoch
	latestAtt,err := protector.RetrieveLatestAttestation(account)
	if err != nil {
		return nil,err
	}
	if latestAtt != nil {
		lookupEndEpoch = latestAtt.Target.Epoch
	}

	history,err := protector.store.ListAttestations(account, lookupStartEpoch, lookupEndEpoch)
	if err != nil {
		return nil,err
	}

	return data.SlashesAttestations(history), nil
}

func (protector *NormalProtection) IsSlashableProposal(account Account, req *pb.SignBeaconProposalRequest) (*ProposalSlashStatus,error) {
	matchedProposal,err := protector.store.RetrieveProposal(account,req.Data.Slot)
	if err != nil && err.Error() != "proposal not found" {
		return nil, err
	}

	if matchedProposal == nil {
		return nil,nil
	}

	data := ToCoreBlockData(req)

	// if it's the same
	if data.Compare(matchedProposal) {
		return nil, nil
	}

	// slashable
	return &ProposalSlashStatus{
		Proposal: data,
		Status:   DoubleProposal,
	},nil
}

func (protector *NormalProtection) SaveAttestation(account Account, req *pb.SignBeaconAttestationRequest) error {
	data := ToCoreAttestationData(req)
	err := protector.store.SaveAttestation(account,data)
	if err != nil {
		return err
	}
	return protector.SaveLatestAttestation(account,req)
}

func (protector *NormalProtection) SaveProposal(account Account, req *pb.SignBeaconProposalRequest) error {
	data := ToCoreBlockData(req)
	return protector.store.SaveProposal(account,data)
}

func (protector *NormalProtection) SaveLatestAttestation(account Account, req *pb.SignBeaconAttestationRequest) error {
	val,err := protector.store.RetrieveLatestAttestation(account)
	if err != nil {
		return nil
	}

	data := ToCoreAttestationData(req)
	if val == nil {
		return protector.store.SaveLatestAttestation(account,data)
	}
	if val.Target.Epoch < req.Data.Target.Epoch { // only write newer
		return protector.store.SaveLatestAttestation(account,data)
	}

	return nil
}

func (protector *NormalProtection) RetrieveLatestAttestation(account Account) (*BeaconAttestation, error) {
	return protector.store.RetrieveLatestAttestation(account)
}

// specialized func that will prevent overflow for lookup epochs for uint64
func lookupEpochSub(l uint64, r uint64) uint64 {
	if l >= r {
		return l-r
	}
	return 0
}