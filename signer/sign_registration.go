package signer

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/signing"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// SignRegistration signs the given ValidatorRegistration.
func (signer *SimpleSigner) SignRegistration(registration *ethpb.ValidatorRegistrationV1, domain []byte, pubKey []byte) ([]byte, error) {
	// Validate the registration.
	if registration == nil {
		return nil, errors.New("registration data is nil")
	}

	// Get the account.
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}
	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	// Produce the signature.
	root, err := signing.ComputeSigningRoot(registration, domain)
	if err != nil {
		return nil, err
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}
	return sig, nil
}
