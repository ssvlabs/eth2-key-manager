package ValidatorSigner

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

type VaultSlashingProtector interface {
	IsSlashableAttestation(req *pb.SignBeaconAttestationRequest) (bool,error)
	IsSlashablePropose(req *pb.SignBeaconProposalRequest) (bool,error)
	SaveAttestation(req *pb.SignBeaconAttestationRequest) error
	SaveProposal(req *pb.SignBeaconProposalRequest) error
}
