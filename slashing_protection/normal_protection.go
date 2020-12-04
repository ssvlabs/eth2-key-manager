package slashingprotection

import (
	"fmt"

	"github.com/pkg/errors"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
)

// NormalProtection implements normal protection logic
type NormalProtection struct {
	store core.SlashingStore
}

// NewNormalProtection is the constructor of NormalProtection
func NewNormalProtection(store core.SlashingStore) *NormalProtection {
	return &NormalProtection{store: store}
}

// IsSlashableAttestation detects double, surround and surrounded slashable events
func (protector *NormalProtection) IsSlashableAttestation(pubKey []byte, attestation *eth.AttestationData) (*core.AttestationSlashStatus, error) {
	// lookupEndEpoch should be the latest written attestation, if not than req.Data.Target.Epoch
	highest, err := protector.RetrieveHighestAttestation(pubKey)
	if err != nil {
		return nil, err
	}
	if highest != nil {
		// Source epoch can't be lower than previously known highest source, it can be equal or higher.
		if attestation.Source.Epoch < highest.Source.Epoch || attestation.Target.Epoch <= highest.Target.Epoch {
			return &core.AttestationSlashStatus{
				Attestation: attestation,
				Status:      core.HighestAttestationVote,
			}, nil
		}
		return nil, nil
	}

	return nil, errors.New("highest attestation data is nil, can't determine if attestation is slashable")
}

// IsSlashableProposal detects slashable proposal request
func (protector *NormalProtection) IsSlashableProposal(pubKey []byte, block *eth.BeaconBlock) (*core.ProposalSlashStatus, error) {
	highest := protector.store.RetrieveHighestProposal(pubKey)
	if highest == nil {
		return nil, fmt.Errorf("highest proposal data is nil, can't determine if proposal is slashable")
	}

	if block.Slot > highest.Slot {
		return &core.ProposalSlashStatus{
			Proposal: nil,
			Status:   core.ValidProposal,
		}, nil
	}

	return &core.ProposalSlashStatus{
		Proposal: nil,
		Status:   core.HighestProposalVote,
	}, nil
}

// UpdateHighestAttestation potentially updates the highest attestation given this latest attestation.
func (protector *NormalProtection) UpdateHighestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
	// if no previous highest attestation found, set current
	highest := protector.store.RetrieveHighestAttestation(pubKey)
	if highest == nil {
		err := protector.store.SaveHighestAttestation(pubKey, attestation)
		if err != nil {
			return err
		}
		return nil
	}

	// Taken from https://github.com/prysmaticlabs/prysm/blob/master/slasher/detection/detect.go#L233
	shouldUpdate := false
	if highest.Source.Epoch < attestation.Source.Epoch {
		highest.Source.Epoch = attestation.Source.Epoch
		shouldUpdate = true
	}
	if highest.Target.Epoch < attestation.Target.Epoch {
		highest.Target.Epoch = attestation.Target.Epoch
		shouldUpdate = true
	}

	if shouldUpdate {
		err := protector.store.SaveHighestAttestation(pubKey, highest)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateHighestProposal updates highest proposal
func (protector *NormalProtection) UpdateHighestProposal(key []byte, block *eth.BeaconBlock) error {
	// if no previous highest proposal found, set current
	highest := protector.store.RetrieveHighestProposal(key)
	if highest == nil {
		err := protector.store.SaveHighestProposal(key, block)
		if err != nil {
			return err
		}
		return nil
	}

	if highest.Slot < block.Slot {
		return protector.store.SaveHighestProposal(key, block)
	}

	return nil
}

// RetrieveHighestAttestation returns highest attestation data
func (protector *NormalProtection) RetrieveHighestAttestation(pubKey []byte) (*eth.AttestationData, error) {
	return protector.store.RetrieveHighestAttestation(pubKey), nil
}
