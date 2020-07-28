package core

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type SlashingProtector interface {
	IsSlashableAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) ([]*AttestationSlashStatus,error)
	IsSlashableProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) *ProposalSlashStatus
	SaveAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error
	SaveProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) error
	SaveLatestAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error
	RetrieveLatestAttestation(key e2types.PublicKey) (*BeaconAttestation, error)
}

type SlashingStore interface {
	SaveAttestation(key e2types.PublicKey, req *BeaconAttestation) error
	RetrieveAttestation(key e2types.PublicKey, epoch uint64) (*BeaconAttestation, error)
	// both epochStart and epochEnd reflect saved attestations by their target epoch
	ListAttestations(key e2types.PublicKey, epochStart uint64, epochEnd uint64) ([]*BeaconAttestation, error)
	SaveProposal(key e2types.PublicKey, req *BeaconBlockHeader) error
	RetrieveProposal(key e2types.PublicKey, slot uint64) (*BeaconBlockHeader, error)
	SaveLatestAttestation(key e2types.PublicKey, req *BeaconAttestation) error
	RetrieveLatestAttestation(key e2types.PublicKey) (*BeaconAttestation, error)
}
