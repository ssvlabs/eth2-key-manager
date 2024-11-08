package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/eth2-key-manager/core"
)

// SignBlock signs the given beacon block
func (signer *SimpleSigner) SignBlock(block ssz.HashRoot, slot phase0.Slot, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	// 2. lock for current account
	val := signer.lock(account.ID(), "proposal")
	val.Lock()
	defer val.Unlock()

	// 3. far future check
	if !IsValidFarFutureSlot(signer.network, slot) {
		return nil, nil, errors.Errorf("proposed block slot too far into the future")
	}

	// 4. check we can even sign this
	status, err := signer.slashingProtector.IsSlashableProposal(pubKey, slot)
	if err != nil {
		return nil, nil, err
	}
	if status.Status != core.ValidProposal {
		return nil, nil, errors.Errorf("slashable proposal (%s), not signing", status.Status)
	}

	// 5. add to protection storage
	if err = signer.slashingProtector.UpdateHighestProposal(pubKey, slot); err != nil {
		return nil, nil, err
	}

	root, err := ComputeETHSigningRoot(block, domain)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}
