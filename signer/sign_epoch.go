package signer

import (
	"encoding/hex"

	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/signing"
)

// SignEpoch signs the given epoch
func (signer *SimpleSigner) SignEpoch(epoch types.Epoch, domain []byte, pubKey []byte) ([]byte, error) {
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

	root, err := signing.ComputeSigningRoot(epoch, domain)
	if err != nil {
		return nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
