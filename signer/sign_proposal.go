package signer

import (
	"encoding/hex"

	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/prysmaticlabs/prysm/runtime/version"

	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/block"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/signing"

	"github.com/bloxapp/eth2-key-manager/core"
)

// TEMPORARYPhase0BlockConversion takes prysm's beacon block interface and converts it to phase0 block
func TEMPORARYPhase0BlockConversion(b block.BeaconBlock) *ethpb.BeaconBlock {
	return &ethpb.BeaconBlock{
		ProposerIndex: b.ProposerIndex(),
		Slot:          b.Slot(),
		ParentRoot:    b.ParentRoot(),
		StateRoot:     b.StateRoot(),
		Body: &ethpb.BeaconBlockBody{
			RandaoReveal:      b.Body().RandaoReveal(),
			Eth1Data:          b.Body().Eth1Data(),
			Graffiti:          b.Body().Graffiti(),
			ProposerSlashings: b.Body().ProposerSlashings(),
			AttesterSlashings: b.Body().AttesterSlashings(),
			Attestations:      b.Body().Attestations(),
			Deposits:          b.Body().Deposits(),
			VoluntaryExits:    b.Body().VoluntaryExits(),
		},
	}
}

// SignBeaconBlock signs the given beacon block
func (signer *SimpleSigner) SignBeaconBlock(b block.BeaconBlock, domain []byte, pubKey []byte) ([]byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "proposal")
	defer signer.unlock(account.ID(), "proposal")

	// 3. far future check
	if !IsValidFarFutureSlot(signer.network, b.Slot()) {
		return nil, errors.Errorf("proposed block slot too far into the future")
	}

	// 4. check we can even sign this
	status, err := signer.verifySlashableAndUpdate(b, pubKey)
	if err != nil {
		return nil, err
	}
	if status.Status != core.ValidProposal {
		return nil, errors.Errorf("slashable proposal (%s), not signing", status.Status)
	}

	// 5. generate ssz root hash and sign
	var root [32]byte
	switch b.Version() {
	case version.Bellatrix:
		block, ok := b.Proto().(*ethpb.BeaconBlockBellatrix)
		if !ok {
			return nil, errors.New("failed type assertion for bellatrix block")
		}
		root, err = signing.ComputeSigningRoot(block, domain)
		if err != nil {
			return nil, errors.Wrap(err, "could not get signing root")
		}
	case version.Altair:
		block, ok := b.Proto().(*ethpb.BeaconBlockAltair)
		if !ok {
			return nil, errors.New("failed type assertion for altair block")
		}
		root, err = signing.ComputeSigningRoot(block, domain)
		if err != nil {
			return nil, errors.Wrap(err, "could not get signing root")
		}
	case version.Phase0:
		block, ok := b.Proto().(*ethpb.BeaconBlock)
		if !ok {
			return nil, errors.New("failed type assertion for phase0 block")
		}
		root, err = signing.ComputeSigningRoot(block, domain)
		if err != nil {
			return nil, errors.Wrap(err, "could not get signing root")
		}
	default:
		return nil, errors.Errorf("unsupported block version %d", b.Version())
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// verifySlashableAndUpdate verified if block is slashable, if not saves it as the highest
func (signer *SimpleSigner) verifySlashableAndUpdate(b block.BeaconBlock, pubKey []byte) (*core.ProposalSlashStatus, error) {
	/**
	We convert the beacon block interface into a phase 0 block, we can allow to do so (even with the differences between phase0 and altair blocks)
	because slashing conditions didn't change.
	TODO - clean up clear separation between phase0 and altair
	*/
	phase0Blk := TEMPORARYPhase0BlockConversion(b)
	status, err := signer.slashingProtector.IsSlashableProposal(pubKey, phase0Blk)
	if err != nil {
		return nil, err
	}

	if err := signer.slashingProtector.UpdateHighestProposal(pubKey, phase0Blk); err != nil {
		return nil, err
	}
	return status, nil
}
