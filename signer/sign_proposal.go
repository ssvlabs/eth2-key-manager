package signer

import (
	"encoding/hex"

	"github.com/pkg/errors"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"

	"github.com/bloxapp/eth2-key-manager/core"
)

// SignBeaconBlock signs the given beacon block
func (signer *SimpleSigner) SignBeaconBlock(block *eth.BeaconBlock, domain []byte, pubKey []byte) ([]byte, error) {
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
	if !IsValidFarFutureSlot(signer.network, block.Slot) {
		return nil, errors.Errorf("proposed block slot too far into the future")
	}

	// 4. check we can even sign this
	status, err := signer.slashingProtector.IsSlashableProposal(pubKey, block)
	if err != nil {
		return nil, err
	}
	if status.Status != core.ValidProposal {
		return nil, errors.Errorf("slashable proposal (%s), not signing", status.Status)
	}

	// 5. add to protection storage
	if err := signer.slashingProtector.UpdateHighestProposal(pubKey, block); err != nil {
		return nil, err
	}

	// 6. generate ssz root hash and sign
	root, err := helpers.ComputeSigningRoot(block, domain)
	if err != nil {
		return nil, err
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
