package in_memory

import (
	"encoding/hex"
	"fmt"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/pkg/errors"
)

func (store *InMemStore) SaveHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	store.highestAttestation[hex.EncodeToString(pubKey)] = attestation
	return nil
}

func (store *InMemStore) RetrieveHighestAttestation(pubKey []byte) *eth.AttestationData {
	if val, ok := store.highestAttestation[hex.EncodeToString(pubKey)]; ok {
		return val
	}
	return nil
}

func (store *InMemStore) SaveProposal(pubKey []byte, block *eth.BeaconBlock) error {
	store.proposalMemory[proposalKey(pubKey, block.Slot)] = block
	return nil
}

func (store *InMemStore) RetrieveProposal(pubKey []byte, slot uint64) (*eth.BeaconBlock, error) {
	ret := store.proposalMemory[proposalKey(pubKey, slot)]
	if ret == nil {
		return nil, errors.New("proposal not found")
	}
	return ret, nil
}

func proposalKey(pubKey []byte, targetSlot uint64) string {
	return fmt.Sprintf("%s_%d", hex.EncodeToString(pubKey), targetSlot)
}
