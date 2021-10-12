package signer

import (
	"encoding/hex"

	"github.com/pkg/errors"

	"github.com/prysmaticlabs/prysm/beacon-chain/core/signing"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// SignAggregateAndProof signs aggregate and proof
func (signer *SimpleSigner) SignAggregateAndProof(agg *eth.AggregateAttestationAndProof, domain []byte, pubKey []byte) ([]byte, error) {
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

	root, err := signing.ComputeSigningRoot(agg, domain)
	if err != nil {
		return nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
