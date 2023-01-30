package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"

	"github.com/pkg/errors"
)

// SignSlot signes the given slot
func (signer *SimpleSigner) SignSlot(slot phase0.Slot, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. check we can even sign this
	// TODO - should we?

	// 2. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	root, err := types.ComputeETHSigningRoot(types.SSZUint64(slot), domain)
	if err != nil {
		return nil, nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}
