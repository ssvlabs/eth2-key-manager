package signer

import (
	"encoding/hex"

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// SignRegistration signs the given ValidatorRegistration.
func (signer *SimpleSigner) SignRegistration(registration *apiv1.ValidatorRegistration, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// Validate the registration.
	if registration == nil {
		return nil, nil, errors.New("registration data is nil")
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
	root, err := types.ComputeETHSigningRoot(registration, domain)
	if err != nil {
		return nil, nil, err
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}
	return sig, root[:], nil
}
