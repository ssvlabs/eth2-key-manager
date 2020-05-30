package in_memory

import (
	"fmt"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func (store *InMemStore) SaveAttestation(account types.Account, req *pb.SignBeaconAttestationRequest) error {
	store.attMemory[attestationKey(account,req.Data.Target.Epoch)] = req
	return nil
}

func (store *InMemStore) RetrieveAttestation(account types.Account, epoch uint64, slot uint64) (*pb.SignBeaconAttestationRequest, error) {
	ret := store.attMemory[attestationKey(account,epoch)]
	if ret == nil {
		return nil,fmt.Errorf("attestation not found")
	}
	return ret,nil
}

func (store *InMemStore) SaveProposal(account types.Account, req *pb.SignBeaconProposalRequest) error {
	store.proposalMemory[proposalKey(account,req.Data.Slot)] = req
	return nil
}

func (store *InMemStore) RetrieveProposal(account types.Account, epoch uint64, slot uint64) (*pb.SignBeaconProposalRequest, error) {
	ret := store.proposalMemory[proposalKey(account,epoch)]
	if ret == nil {
		return nil,fmt.Errorf("proposal not found")
	}
	return ret,nil
}

func attestationKey(account types.Account, targetEpoch uint64) string {
	return fmt.Sprintf("%s_%d",account.ID().String(),targetEpoch)
}

func proposalKey(account types.Account, targetSlot uint64) string {
	return fmt.Sprintf("%s_%d",account.ID().String(),targetSlot)
}