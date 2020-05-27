package slashing_protectors

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type VaultSlashingProtector interface {
	IsSlashableAttestation(req *pb.SignBeaconAttestationRequest) (bool,error)
	IsSlashablePropose(req *pb.SignBeaconProposalRequest) (bool,error)
	SaveAttestation(req *pb.SignBeaconAttestationRequest) error
	SaveProposal(req *pb.SignBeaconProposalRequest) error
}

type NormalProtection struct {
	store types.Store
}

func NewNormalProtection(store types.Store) *NormalProtection {
	return &NormalProtection{store:store}
}

func (protector *NormalProtection) IsSlashableAttestation(req *pb.SignBeaconAttestationRequest) (bool,error) {
	return false,nil
}

func (protector *NormalProtection) IsSlashablePropose(req *pb.SignBeaconProposalRequest) (bool,error) {
	return false,nil
}

func (protector *NormalProtection) SaveAttestation(req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (protector *NormalProtection) SaveProposal(req *pb.SignBeaconProposalRequest) error {
	return nil
}