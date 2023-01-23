package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/bloxapp/eth2-key-manager/core"
)

// SignBeaconBlock signs the given beacon block
func (signer *SimpleSigner) SignBeaconBlock(b *spec.VersionedBeaconBlock, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "proposal")
	defer signer.unlock(account.ID(), "proposal")

	slot, err := b.Slot()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get block slot")
	}

	// 3. far future check
	if !IsValidFarFutureSlot(signer.network, slot) {
		return nil, nil, errors.Errorf("proposed block slot too far into the future")
	}

	// 4. check we can even sign this
	status, err := signer.verifySlashableAndUpdate(pubKey, slot)
	if err != nil {
		return nil, nil, err
	}
	if status.Status != core.ValidProposal {
		return nil, nil, errors.Errorf("slashable proposal (%s), not signing", status.Status)
	}

	var block ssz.HashRoot
	switch b.Version {
	case spec.DataVersionPhase0:
		block = b.Phase0
	case spec.DataVersionAltair:
		block = b.Altair
	case spec.DataVersionBellatrix:
		block = b.Bellatrix
	case spec.DataVersionCapella:
		block = b.Capella
	default:
		return nil, nil, errors.Errorf("unsupported block version %d", b.Version)
	}

	root, err := types.ComputeETHSigningRoot(block, domain)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}

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
