package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

// SignAggregateAndProof signs aggregate and proof.
// It can be *phase0.AggregateAndProof or *electra.AggregateAndProof since electra.
// As we don't use any AggregateAndProof's fields, we can just use ssz.HashRoot.
func (signer *SimpleSigner) SignAggregateAndProof(agg ssz.HashRoot, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
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

	root, err := ComputeETHSigningRoot(agg, domain)
	if err != nil {
		return nil, nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}
