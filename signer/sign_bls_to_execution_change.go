package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SignBLSToExecutionChange signs the given BLSToExecutionChange. OFFLINE operation
func (signer *SimpleSigner) SignBLSToExecutionChange(blsToExecutionChange *capella.BLSToExecutionChange, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// Validate the bls to execution change.
	if blsToExecutionChange == nil {
		return nil, nil, errors.New("bls to execution change is nil")
	}

	// Get the account.
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}
	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	// Produce the signature.
	root, err := ComputeETHSigningRoot(blsToExecutionChange, domain)
	if err != nil {
		return nil, nil, err
	}

	// This is actually withdrawal key
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}
	return sig, root[:], nil
}
