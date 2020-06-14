package core

import pb "github.com/wealdtech/eth2-signer-api/pb/v1"

type SlashingProtector interface {
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
