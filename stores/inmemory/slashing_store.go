package inmemory

import (
	"encoding/hex"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// SaveHighestAttestation saves the given highest attestation
func (store *InMemStore) SaveHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	store.highestAttestationLock.Lock()
	store.highestAttestation[hex.EncodeToString(pubKey)] = attestation
	store.highestAttestationLock.Unlock()
	return nil
}

// RetrieveHighestAttestation retrieves highest attestation
func (store *InMemStore) RetrieveHighestAttestation(pubKey []byte) *eth.AttestationData {
	store.highestAttestationLock.RLock()
	val := store.highestAttestation[hex.EncodeToString(pubKey)]
	store.highestAttestationLock.RUnlock()
	return val
}

// SaveHighestProposal saves the given highest attestation
func (store *InMemStore) SaveHighestProposal(pubKey []byte, block *eth.BeaconBlock) error {
	store.highestProposalLock.Lock()
	store.highestProposal[hex.EncodeToString(pubKey)] = block
	store.highestProposalLock.Unlock()
	return nil
}

// RetrieveHighestProposal returns highest proposal
func (store *InMemStore) RetrieveHighestProposal(pubKey []byte) *eth.BeaconBlock {
	store.highestProposalLock.RLock()
	val := store.highestProposal[hex.EncodeToString(pubKey)]
	store.highestProposalLock.RUnlock()
	return val
}
