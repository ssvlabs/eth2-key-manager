package core

import (
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// SlashingProtector represents the behavior of the slashing protector
type SlashingProtector interface {
	IsSlashableAttestation(pubKey []byte, attestation *eth.AttestationData) (*AttestationSlashStatus, error)
	IsSlashableProposal(pubKey []byte, block *eth.BeaconBlock) (*ProposalSlashStatus, error)
	// Will potentially update the highest attestation given this latest attestation.
	UpdateHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error
	UpdateHighestProposal(pubKey []byte, block *eth.BeaconBlock) error
	RetrieveHighestAttestation(pubKey []byte) (*eth.AttestationData, error)
}

// SlashingStore represents the behavior of the slashing store
type SlashingStore interface {
	SaveHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error
	RetrieveHighestAttestation(pubKey []byte) *eth.AttestationData
	SaveHighestProposal(pubKey []byte, block *eth.BeaconBlock) error
	RetrieveHighestProposal(pubKey []byte) *eth.BeaconBlock
}
