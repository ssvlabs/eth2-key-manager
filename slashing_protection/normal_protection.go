package slashing_protection

import (
	"bytes"
	"fmt"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
)

type NormalProtection struct {
	store core.SlashingStore
}

// NewNormalProtection is the constructor of NormalProtection
func NewNormalProtection(store core.SlashingStore) *NormalProtection {
	return &NormalProtection{store: store}
}

// will detect double, surround and surrounded slashable events
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
	} else {
		return nil, fmt.Errorf("highest attestation data is nil, can't determine if attestation is slashable")
	}
}

func (protector *NormalProtection) IsSlashableProposal(pubKey []byte, block *eth.BeaconBlock) *core.ProposalSlashStatus {
	matchedProposal, err := protector.store.RetrieveProposal(pubKey, block.Slot)
	if err != nil && err.Error() != "proposal not found" {
		return &core.ProposalSlashStatus{
			Proposal: nil,
			Status:   core.Error,
			Error:    err,
		}
	}

	if matchedProposal == nil {
		return &core.ProposalSlashStatus{
			Proposal: nil,
			Status:   core.ValidProposal,
		}
	}

	equal := func(a *eth.BeaconBlock, b *eth.BeaconBlock) bool {
		aRoot, err := a.HashTreeRoot()
		if err != nil {
			return false
		}
		bRoot, err := b.HashTreeRoot()
		if err != nil {
			return false
		}
		return bytes.Equal(aRoot[:], bRoot[:])
	}

	// if it's the same
	if equal(block, matchedProposal) {
		return &core.ProposalSlashStatus{
			Proposal: matchedProposal,
			Status:   core.ValidProposal,
		}
	}

	// slashable
	return &core.ProposalSlashStatus{
		Proposal: matchedProposal,
		Status:   core.DoubleProposal,
	}
}

// Will potentially update the highest attestation given this latest attestation.
func (protector *NormalProtection) UpdateLatestAttestation(pubKey []byte, attestation *eth.AttestationData) error {
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

func (protector *NormalProtection) SaveProposal(key []byte, block *eth.BeaconBlock) error {
	return protector.store.SaveProposal(key, block)
}

func (protector *NormalProtection) RetrieveHighestAttestation(pubKey []byte) (*eth.AttestationData, error) {
	return protector.store.RetrieveHighestAttestation(pubKey), nil
}
