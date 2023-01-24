package signer

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/bloxapp/eth2-key-manager/core"
)

// verifySlashableAndUpdate verified if block is slashable, if not saves it as the highest
func (signer *SimpleSigner) verifySlashableAndUpdate(pubKey []byte, slot phase0.Slot) (*core.ProposalSlashStatus, error) {
	status, err := signer.slashingProtector.IsSlashableProposal(pubKey, slot)
	if err != nil {
		return nil, err
	}

	if err := signer.slashingProtector.UpdateHighestProposal(pubKey, slot); err != nil {
		return nil, err
	}
	return status, nil
}
