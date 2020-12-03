package in_memory

import (
	"encoding/hex"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
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

func (store *InMemStore) SaveHighestProposal(pubKey []byte, block *eth.BeaconBlock) error {
	store.highestProposal[hex.EncodeToString(pubKey)] = block
	return nil
}

func (store *InMemStore) RetrieveHighestProposal(pubKey []byte) *eth.BeaconBlock {
	return store.highestProposal[hex.EncodeToString(pubKey)]
}
