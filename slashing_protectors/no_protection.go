package slashing_protectors
import (
	"github.com/bloxapp/KeyVault/core"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type NoProtection struct {

}

func (p *NoProtection)IsSlashableAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) ([]*core.AttestationSlashStatus,error) {
	return make([]*core.AttestationSlashStatus,0),nil
}

func (p *NoProtection)IsSlashableProposal(account types.Account, req *pb.SignBeaconProposalRequest) (*core.ProposalSlashStatus,error) {
	return nil,nil
}

func (p *NoProtection)SaveAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (p *NoProtection)SaveProposal(account types.Account, req *pb.SignBeaconProposalRequest) error {
	return nil
}

func (p *NoProtection)SaveLatestAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) error {
	return nil
}

func (p *NoProtection)RetrieveLatestAttestation(account types.Account) (*core.BeaconAttestation, error) {
	return nil,nil
}
