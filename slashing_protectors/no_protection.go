package slashing_protectors
import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type NoProtection struct {
	store types.Store
}

func (protector *NoProtection) IsSlashableAttestation(req *pb.SignBeaconAttestationRequest) (bool,error) {
	return false,nil
}

func (protector *NoProtection) IsSlashablePropose(req *pb.SignBeaconProposalRequest) (bool,error) {
	return false,nil
}

func (protector *NoProtection) SaveAttestation(req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (protector *NoProtection) SaveProposal(req *pb.SignBeaconProposalRequest) error {
	return nil
}
