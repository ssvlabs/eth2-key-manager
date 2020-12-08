package inmemory

import (
	"encoding/hex"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SaveHighestAttestation saves the given highest attestation
func (store *InMemStore) SaveHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	store.highestAttestation[hex.EncodeToString(pubKey)] = attestation
	return nil
}

// RetrieveHighestAttestation retrieves highest attestation
func (store *InMemStore) RetrieveHighestAttestation(pubKey []byte) *eth.AttestationData {
	if val, ok := store.highestAttestation[hex.EncodeToString(pubKey)]; ok {
		return val
	}
	return nil
}

// SaveHighestProposal saves the given highest attestation
func (store *InMemStore) SaveHighestProposal(pubKey []byte, block *eth.BeaconBlock) error {
	store.highestProposal[hex.EncodeToString(pubKey)] = block
	return nil
}

// RetrieveHighestProposal returns highest proposal
func (store *InMemStore) RetrieveHighestProposal(pubKey []byte) *eth.BeaconBlock {
	return store.highestProposal[hex.EncodeToString(pubKey)]
}
