package in_memory

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

func (store *InMemStore) SaveAttestation(key e2types.PublicKey, req *core.BeaconAttestation) error {
	store.attMemory[attestationKey(key, req.Target.Epoch)] = req
	return nil
}

func (store *InMemStore) RetrieveAttestation(key e2types.PublicKey, epoch uint64) (*core.BeaconAttestation, error) {
	ret := store.attMemory[attestationKey(key, epoch)]
	if ret == nil {
		return nil, fmt.Errorf("attestation not found")
	}
	return ret, nil
}

func (store *InMemStore) ListAttestations(key e2types.PublicKey, epochStart uint64, epochEnd uint64) ([]*core.BeaconAttestation, error) {
	ret := make([]*core.BeaconAttestation, 0)
	for i := epochStart; i <= epochEnd; i++ {
		if val, err := store.RetrieveAttestation(key, i); val != nil && err == nil {
			ret = append(ret, val)
		}
	}
	return ret, nil
}

func (store *InMemStore) SaveProposal(key e2types.PublicKey, req *core.BeaconBlockHeader) error {
	store.proposalMemory[proposalKey(key, req.Slot)] = req
	return nil
}

func (store *InMemStore) RetrieveProposal(key e2types.PublicKey, slot uint64) (*core.BeaconBlockHeader, error) {
	ret := store.proposalMemory[proposalKey(key, slot)]
	if ret == nil {
		return nil, fmt.Errorf("proposal not found")
	}
	return ret, nil
}

func (store *InMemStore) SaveLatestAttestation(key e2types.PublicKey, req *core.BeaconAttestation) error {
	store.attMemory[hex.EncodeToString(key.Marshal())+"_latest"] = req
	return nil
}

func (store *InMemStore) RetrieveLatestAttestation(key e2types.PublicKey) (*core.BeaconAttestation, error) {
	return store.attMemory[hex.EncodeToString(key.Marshal())+"_latest"], nil
}

func attestationKey(key e2types.PublicKey, targetEpoch uint64) string {
	return fmt.Sprintf("%s_%d", hex.EncodeToString(key.Marshal()), targetEpoch)
}

func proposalKey(key e2types.PublicKey, targetSlot uint64) string {
	return fmt.Sprintf("%s_%d", hex.EncodeToString(key.Marshal()), targetSlot)
}
