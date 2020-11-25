package in_memory

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

func (store *InMemStore) SaveHighestAttestation(key e2types.PublicKey, req *core.BeaconAttestation) error {
	store.highestAttestation = req
	return nil
}

func (store *InMemStore) RetrieveHighestAttestation(key e2types.PublicKey) *core.BeaconAttestation {
	return store.highestAttestation
}

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

func attestationKey(key e2types.PublicKey) string {
	return hex.EncodeToString(key.Marshal())
}

func proposalKey(key e2types.PublicKey, targetSlot uint64) string {
	return fmt.Sprintf("%s_%d", hex.EncodeToString(key.Marshal()), targetSlot)
}
