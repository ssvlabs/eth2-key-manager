package slashing_protection

import (
	"github.com/bloxapp/KeyVault/core"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type NoProtection struct {

}

func (p *NoProtection)IsSlashableAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) ([]*core.AttestationSlashStatus,error) {
	return make([]*core.AttestationSlashStatus,0),nil
}

func (p *NoProtection)IsSlashableProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) *core.ProposalSlashStatus {
	return &core.ProposalSlashStatus{
		Proposal: nil,
		Status:   core.ValidProposal,
	}
}

func (p *NoProtection)SaveAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (p *NoProtection)SaveProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) error {
	return nil
}

func (p *NoProtection)SaveLatestAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (p *NoProtection)RetrieveLatestAttestation(key e2types.PublicKey) (*core.BeaconAttestation, error) {
	return nil,nil
}
