package slashing_protection

import (
	"github.com/bloxapp/eth2-key-manager/core"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

type NoProtection struct {
}

func (p *NoProtection) IsSlashableAttestation(pubKey []byte, attestation *eth.AttestationData) (*core.AttestationSlashStatus, error) {
	return nil, nil
}

func (p *NoProtection) IsSlashableProposal(pubKey []byte, block *eth.BeaconBlock) *core.ProposalSlashStatus {
	return &core.ProposalSlashStatus{
		Proposal: nil,
		Status:   core.ValidProposal,
	}
}

func (p *NoProtection) SaveProposal(pubKey []byte, block *eth.BeaconBlock) error {
	return nil
}

func (p *NoProtection) UpdateLatestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	return nil
}

func (p *NoProtection) RetrieveHighestAttestation(pubKey []byte) (*eth.AttestationData, error) {
	return nil, nil
}
