package slashing_protectors

import (
	"github.com/bloxapp/KeyVault/core"
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
	SaveAttestation(account types.Account, req *core.BeaconAttestation) error
	RetrieveAttestation(account types.Account, epoch uint64) (*core.BeaconAttestation, error)
	ListAttestations(account types.Account, epochStart uint64, epochEnd uint64) ([]*core.BeaconAttestation, error)
	SaveProposal(account types.Account, req *core.BeaconBlockHeader) error
	RetrieveProposal(account types.Account, epoch uint64) (*core.BeaconBlockHeader, error)
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
	data := core.ToCoreAttestationData(req)
	return protector.store.SaveAttestation(account,data)
}

func (protector *NormalProtection) SaveProposal(account types.Account, req *pb.SignBeaconProposalRequest) error {
	data := core.ToCoreBlockData(req)
	return protector.store.SaveProposal(account,data)
}