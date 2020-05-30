package slashing_protectors

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type VaultSlashingProtector interface {
	IsSlashableAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) (bool,error)
	IsSlashablePropose(account types.Account, req *pb.SignBeaconProposalRequest) (bool,error)
	SaveAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) error
	SaveProposal(account types.Account, req *pb.SignBeaconProposalRequest) error
}

type SlashingStore interface {
	SaveAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) error
	RetrieveAttestation(account types.Account, epoch uint64, slot uint64) (*pb.SignBeaconAttestationRequest, error)
	SaveProposal(account types.Account, req *pb.SignBeaconProposalRequest) error
	RetrieveProposal(account types.Account, epoch uint64, slot uint64) (*pb.SignBeaconProposalRequest, error)
}

type NormalProtection struct {
	store SlashingStore
}

func NewNormalProtection(store SlashingStore) *NormalProtection {
	return &NormalProtection{store:store}
}

func (protector *NormalProtection) IsSlashableAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) (bool,error) {
	return false,nil
}

func (protector *NormalProtection) IsSlashablePropose(account types.Account, req *pb.SignBeaconProposalRequest) (bool,error) {
	return false,nil
}

func (protector *NormalProtection) SaveAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) error {
	return protector.store.SaveAttestation(account,req)
}

func (protector *NormalProtection) SaveProposal(account types.Account, req *pb.SignBeaconProposalRequest) error {
	return nil
}