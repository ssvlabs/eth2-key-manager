package core

import (
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SlashingProtector represents the behavior of the slashing protector
type SlashingProtector interface {
	IsSlashableAttestation(pubKey []byte, attestation *eth.AttestationData) (*AttestationSlashStatus, error)
	IsSlashableProposal(pubKey []byte, block *eth.BeaconBlock) *ProposalSlashStatus
	// Will potentially update the highest attestation given this latest attestation.
	UpdateLatestAttestation(pubKey []byte, attestation *eth.AttestationData) error
	SaveProposal(pubKey []byte, block *eth.BeaconBlock) error
	RetrieveHighestAttestation(pubKey []byte) (*eth.AttestationData, error)
}

// SlashingStore represents the behavior of the slashing store
type SlashingStore interface {
	SaveHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error
	RetrieveHighestAttestation(pubKey []byte) *eth.AttestationData
	SaveProposal(pubKey []byte, block *eth.BeaconBlock) error
	RetrieveProposal(pubKey []byte, slot uint64) (*eth.BeaconBlock, error)
}
