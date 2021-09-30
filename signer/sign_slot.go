package signer

import (
	"encoding/hex"

	types "github.com/prysmaticlabs/eth2-types"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
)

// SignSlot signes the given slot
func (signer *SimpleSigner) SignSlot(slot types.Slot, domain []byte, pubKey []byte) ([]byte, error) {
	// 1. check we can even sign this
	// TODO - should we?

	// 2. get the account
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	root, err := helpers.ComputeSigningRoot(slot, domain)
	if err != nil {
		return nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
