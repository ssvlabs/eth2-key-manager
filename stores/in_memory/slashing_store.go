package in_memory

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
)

func (store *InMemStore) SaveAttestation(account core.Account, req *core.BeaconAttestation) error {
	store.attMemory[attestationKey(account,req.Target.Epoch)] = req
	return nil
}

func (store *InMemStore) RetrieveAttestation(account core.Account, epoch uint64) (*core.BeaconAttestation, error) {
	ret := store.attMemory[attestationKey(account,epoch)]
	if ret == nil {
		return nil,fmt.Errorf("attestation not found")
	}
	return ret,nil
}

func (store *InMemStore) ListAttestations(account core.Account, epochStart uint64, epochEnd uint64) ([]*core.BeaconAttestation, error) {
	ret := make([]*core.BeaconAttestation,0)
	for i:= epochStart ; i <= epochEnd ; i++ {
		if val,err := store.RetrieveAttestation(account,i); val != nil && err == nil {
			ret = append(ret, val)
		}
	}
	return ret,nil
}

func (store *InMemStore) SaveProposal(account core.Account, req *core.BeaconBlockHeader) error {
	store.proposalMemory[proposalKey(account,req.Slot)] = req
	return nil
}

func (store *InMemStore) RetrieveProposal(account core.Account, slot uint64) (*core.BeaconBlockHeader, error) {
	ret := store.proposalMemory[proposalKey(account,slot)]
	if ret == nil {
		return nil,fmt.Errorf("proposal not found")
	}
	return ret,nil
}

func (store *InMemStore) SaveLatestAttestation(account core.Account, req *core.BeaconAttestation) error {
	store.attMemory[account.ID().String() + "_latest"] = req
	return nil
}

func (store *InMemStore) RetrieveLatestAttestation(account core.Account) (*core.BeaconAttestation, error) {
	return store.attMemory[account.ID().String() + "_latest"],nil
}

func attestationKey(account core.Account, targetEpoch uint64) string {
	return fmt.Sprintf("%s_%d",account.ID().String(),targetEpoch)
}

func proposalKey(account core.Account, targetSlot uint64) string {
	return fmt.Sprintf("%s_%d",account.ID().String(),targetSlot)
}