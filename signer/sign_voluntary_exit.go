package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/bloxapp/eth2-key-manager/core"
)

// SignVoluntaryExit signs the given VoluntaryExit.
func (signer *SimpleSigner) SignVoluntaryExit(voluntaryExit *phase0.VoluntaryExit, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// Validate the voluntary exit.
	if voluntaryExit == nil {
		return nil, nil, errors.New("voluntary exit data is nil")
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
	root, err := core.ComputeETHSigningRoot(voluntaryExit, domain)
	if err != nil {
		return nil, nil, err
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}
	return sig, root[:], nil
}
