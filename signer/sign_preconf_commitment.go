package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// TODO - add test(s)
// SignSlot signes the given preconf-commitment data
func (signer *SimpleSigner) SignPreconfCommitment(data []byte, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	sszRoot := SSZBytes(data)
	root, err := ComputeETHSigningRoot(sszRoot, domain)
	if err != nil {
		return nil, nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}
