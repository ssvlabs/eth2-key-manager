package core

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// SlashingProtector represents the behavior of the slashing protector
type SlashingProtector interface {
	IsSlashableAttestation(pubKey []byte, attestation *phase0.AttestationData) (*AttestationSlashStatus, error)
	IsSlashableProposal(pubKey []byte, slot phase0.Slot) (*ProposalSlashStatus, error)
	UpdateHighestAttestation(pubKey []byte, attestation *phase0.AttestationData) error
	UpdateHighestProposal(pubKey []byte, slot phase0.Slot) error
	RetrieveHighestAttestation(pubKey []byte) (*phase0.AttestationData, error)
}

// SlashingStore represents the behavior of the slashing store
type SlashingStore interface {
	SaveHighestAttestation(pubKey []byte, attestation *phase0.AttestationData) error
	RetrieveHighestAttestation(pubKey []byte) (*phase0.AttestationData, error)
	SaveHighestProposal(pubKey []byte, slot phase0.Slot) error
	RetrieveHighestProposal(pubKey []byte) (phase0.Slot, bool, error)
}
