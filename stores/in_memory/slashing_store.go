package in_memory

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

func (store *InMemStore) SaveHighestAttestation(key e2types.PublicKey, req *core.BeaconAttestation) error {
	k := fmt.Sprintf("%s_highest", attestationKey(key))
	store.attMemory[k] = req
	return nil
}

func (store *InMemStore) RetrieveHighestAttestation(key e2types.PublicKey) (*core.BeaconAttestation, error) {
	k := fmt.Sprintf("%s_highest", attestationKey(key))
	ret := store.attMemory[k]
	if ret == nil {
		return nil, errors.New("attestation not found")
	}
	return ret, nil
}

//func (store *InMemStore) ListAttestations(key e2types.PublicKey, epochStart uint64, epochEnd uint64) ([]*core.BeaconAttestation, error) {
//	ret := make([]*core.BeaconAttestation, 0)
//	for i := epochStart; i <= epochEnd; i++ {
//		if val, err := store.RetrieveHighestAttestation(key, i); val != nil && err == nil {
//			ret = append(ret, val)
//		}
//	}
//	return ret, nil
//}

func (store *InMemStore) SaveProposal(key e2types.PublicKey, req *core.BeaconBlockHeader) error {
	store.proposalMemory[proposalKey(key, req.Slot)] = req
	return nil
}

func (store *InMemStore) RetrieveProposal(key e2types.PublicKey, slot uint64) (*core.BeaconBlockHeader, error) {
	ret := store.proposalMemory[proposalKey(key, slot)]
	if ret == nil {
		return nil, errors.New("proposal not found")
	}
	return ret, nil
}

//func (store *InMemStore) SaveLatestAttestation(key e2types.PublicKey, req *core.BeaconAttestation) error {
//	store.attMemory[hex.EncodeToString(key.Marshal())+"_latest"] = req
//	return nil
//}
//
//func (store *InMemStore) RetrieveLatestAttestation(key e2types.PublicKey) (*core.BeaconAttestation, error) {
//	return store.attMemory[hex.EncodeToString(key.Marshal())+"_latest"], nil
//}

func attestationKey(key e2types.PublicKey) string {
	return fmt.Sprintf("%s", hex.EncodeToString(key.Marshal()))
}

func proposalKey(key e2types.PublicKey, targetSlot uint64) string {
	return fmt.Sprintf("%s_%d", hex.EncodeToString(key.Marshal()), targetSlot)
}
