package slashing_protection

import (
	"github.com/bloxapp/eth2-key-manager/core"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// NoProtection implements slashing protector interface with dummy implementation
type NoProtection struct {
}

// IsSlashableAttestation returns always nils
func (p *NoProtection) IsSlashableAttestation(pubKey []byte, attestation *eth.AttestationData) (*core.AttestationSlashStatus, error) {
	return nil, nil
}

// IsSlashableProposal returns always valid result
func (p *NoProtection) IsSlashableProposal(pubKey []byte, block *eth.BeaconBlock) (*core.ProposalSlashStatus, error) {
	return &core.ProposalSlashStatus{
		Proposal: nil,
		Status:   core.ValidProposal,
	}, nil
}

// UpdateHighestProposal does nothing
func (p *NoProtection) UpdateHighestProposal(pubKey []byte, block *eth.BeaconBlock) error {
	return nil
}

// UpdateHighestAttestation does nothing
func (p *NoProtection) UpdateHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	return nil
}

// RetrieveHighestAttestation does nothing
func (p *NoProtection) RetrieveHighestAttestation(pubKey []byte) (*eth.AttestationData, error) {
	return nil, nil
}
