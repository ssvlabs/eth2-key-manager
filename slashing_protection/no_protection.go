package slashing_protection

import (
	"github.com/bloxapp/KeyVault/core"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

type NoProtection struct {

}

func (p *NoProtection)IsSlashableAttestation(account core.Account, req *pb.SignBeaconAttestationRequest) ([]*core.AttestationSlashStatus,error) {
	return make([]*core.AttestationSlashStatus,0),nil
}

func (p *NoProtection)IsSlashableProposal(account core.Account, req *pb.SignBeaconProposalRequest) (*core.ProposalSlashStatus,error) {
	return nil,nil
}

func (p *NoProtection)SaveAttestation(account core.Account, req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (p *NoProtection)SaveProposal(account core.Account, req *pb.SignBeaconProposalRequest) error {
	return nil
}

func (p *NoProtection)SaveLatestAttestation(account core.Account, req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (p *NoProtection)RetrieveLatestAttestation(account core.Account) (*core.BeaconAttestation, error) {
	return nil,nil
}
