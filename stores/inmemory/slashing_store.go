package inmemory

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SaveHighestAttestation saves the given highest attestation
func (store *InMemStore) SaveHighestAttestation(pubKey []byte, attestation *phase0.AttestationData) error {
	if pubKey == nil {
		return errors.New("public key could not be nil")
	}

	if attestation == nil {
		return errors.New("attestation data could not be nil")
	}

	store.highestAttestationLock.Lock()
	store.highestAttestation[hex.EncodeToString(pubKey)] = attestation
	store.highestAttestationLock.Unlock()
	return nil
}

// RetrieveHighestAttestation retrieves highest attestation
func (store *InMemStore) RetrieveHighestAttestation(pubKey []byte) (*phase0.AttestationData, bool, error) {
	if pubKey == nil {
		return nil, false, errors.New("public key could not be nil")
	}

	store.highestAttestationLock.RLock()
	val, found := store.highestAttestation[hex.EncodeToString(pubKey)]
	store.highestAttestationLock.RUnlock()
	return val, found, nil
}

// SaveHighestProposal saves the given highest attestation
func (store *InMemStore) SaveHighestProposal(pubKey []byte, slot phase0.Slot) error {
	if pubKey == nil {
		return errors.New("public key could not be nil")
	}
	if slot == 0 {
		return errors.New("invalid proposal slot, slot could not be 0")
	}

	store.highestProposalLock.Lock()
	store.highestProposal[hex.EncodeToString(pubKey)] = uint64(slot)
	store.highestProposalLock.Unlock()
	return nil
}

// RetrieveHighestProposal returns highest proposal
func (store *InMemStore) RetrieveHighestProposal(pubKey []byte) (phase0.Slot, bool, error) {
	if pubKey == nil {
		return 0, false, errors.New("public key could not be nil")
	}

	store.highestProposalLock.RLock()
	val, found := store.highestProposal[hex.EncodeToString(pubKey)]
	store.highestProposalLock.RUnlock()
	return phase0.Slot(val), found, nil
}
