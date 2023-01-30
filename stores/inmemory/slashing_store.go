package inmemory

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// SaveHighestAttestation saves the given highest attestation
func (store *InMemStore) SaveHighestAttestation(pubKey []byte, attestation *phase0.AttestationData) error {
	store.highestAttestationLock.Lock()
	store.highestAttestation[hex.EncodeToString(pubKey)] = attestation
	store.highestAttestationLock.Unlock()
	return nil
}

// RetrieveHighestAttestation retrieves highest attestation
func (store *InMemStore) RetrieveHighestAttestation(pubKey []byte) (*phase0.AttestationData, error) {
	store.highestAttestationLock.RLock()
	val := store.highestAttestation[hex.EncodeToString(pubKey)]
	store.highestAttestationLock.RUnlock()
	return val, nil
}

// SaveHighestProposal saves the given highest attestation
func (store *InMemStore) SaveHighestProposal(pubKey []byte, slot phase0.Slot) error {
	store.highestProposalLock.Lock()
	store.highestProposal[hex.EncodeToString(pubKey)] = slot
	store.highestProposalLock.Unlock()
	return nil
}

// RetrieveHighestProposal returns highest proposal
func (store *InMemStore) RetrieveHighestProposal(pubKey []byte) (phase0.Slot, error) {
	store.highestProposalLock.RLock()
	val := store.highestProposal[hex.EncodeToString(pubKey)]
	store.highestProposalLock.RUnlock()
	return val, nil
}
