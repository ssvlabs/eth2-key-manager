package slashingprotection

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

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
func (protector *NormalProtection) IsSlashableAttestation(pubKey []byte, attestation *phase0.AttestationData) (*core.AttestationSlashStatus, error) {
	// lookupEndEpoch should be the latest written attestation, if not than req.Data.Target.Epoch
	highest, err := protector.RetrieveHighestAttestation(pubKey)
	if err != nil {
		return nil, err
	}
	if highest != nil {
		// Source epoch can't be lower than previously known highest source, it can be equal or higher.
		// We prevent double voting by rejecting another attestations with the same target epoch
		// however you are eligible to sign the message with the same target epoch and the signing root,
		// we are being strict by not storing the signing roots
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
func (protector *NormalProtection) IsSlashableProposal(pubKey []byte, slot phase0.Slot) (*core.ProposalSlashStatus, error) {
	highest, err := protector.store.RetrieveHighestProposal(pubKey)
	if err != nil {
		return nil, errors.New("could not retrieve highest proposal")
	}
	if highest == nil {
		return nil, errors.New("highest proposal data is nil, can't determine if proposal is slashable")
	}

	if slot > *highest {
		return &core.ProposalSlashStatus{
			Slot:   slot,
			Status: core.ValidProposal,
		}, nil
	}

	return &core.ProposalSlashStatus{
		Slot:   slot,
		Status: core.HighestProposalVote,
	}, nil
}

// UpdateHighestAttestation potentially updates the highest attestation given this latest attestation.
func (protector *NormalProtection) UpdateHighestAttestation(pubKey []byte, attestation *phase0.AttestationData) error {
	// if no previous highest attestation found, set current
	highest, err := protector.store.RetrieveHighestAttestation(pubKey)
	if err != nil {
		return errors.New("could not retrieve highest attestation")
	}
	if highest == nil {
		err := protector.store.SaveHighestAttestation(pubKey, attestation)
		if err != nil {
			return errors.New("could not save highest attestation")
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
			return errors.New("could not save highest attestation")
		}
	}
	return nil
}

// UpdateHighestProposal updates highest proposal
func (protector *NormalProtection) UpdateHighestProposal(key []byte, slot phase0.Slot) error {
	// if no previous highest proposal found, set current
	highest, err := protector.store.RetrieveHighestProposal(key)
	if err != nil {
		return errors.New("could not retrieve highest proposal")
	}
	if highest == nil || *highest < slot {
		if err := protector.store.SaveHighestProposal(key, slot); err != nil {
			return errors.New("could not save highest proposal")
		}
	}

	return nil
}

// RetrieveHighestAttestation returns highest attestation data
func (protector *NormalProtection) RetrieveHighestAttestation(pubKey []byte) (*phase0.AttestationData, error) {
	return protector.store.RetrieveHighestAttestation(pubKey)
}
