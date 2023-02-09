package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

// SignRegistration signs the given ValidatorRegistration.
func (signer *SimpleSigner) SignRegistration(registration *api.VersionedValidatorRegistration, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
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

	var reg ssz.HashRoot
	switch registration.Version {
	case spec.BuilderVersionV1:
		if registration.V1 == nil {
			return nil, nil, errors.New("no validator registration")
		}
		reg = registration.V1
	default:
		return nil, nil, errors.Errorf("unsupported registration version %d", registration.Version)
	}

	// Produce the signature.
	root, err := types.ComputeETHSigningRoot(reg, domain)
	if err != nil {
		return nil, nil, err
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}
	return sig, root[:], nil
}
