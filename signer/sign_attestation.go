package signer

import (
	"encoding/hex"

	"github.com/prysmaticlabs/prysm/beacon-chain/core/signing"

	"github.com/pkg/errors"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// SignBeaconAttestation signs beacon attestation data
func (signer *SimpleSigner) SignBeaconAttestation(attestation *eth.AttestationData, domain []byte, pubKey []byte) ([]byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}
	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "attestation")
	defer func() {
		signer.unlock(account.ID(), "attestation")
	}()

	// 3. far future check
	if !IsValidFarFutureEpoch(signer.network, attestation.Target.Epoch) {
		return nil, errors.Errorf("target epoch too far into the future")
	}
	if !IsValidFarFutureEpoch(signer.network, attestation.Source.Epoch) {
		return nil, errors.Errorf("source epoch too far into the future")
	}

	// 4. check we can even sign this
	if val, err := signer.slashingProtector.IsSlashableAttestation(pubKey, attestation); err != nil || val != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.Errorf("slashable attestation (%s), not signing", val.Status)
	}

	// 5. add to protection storage
	if err := signer.slashingProtector.UpdateHighestAttestation(pubKey, attestation); err != nil {
		return nil, err
	}

	// 6. Prepare and sign data
	root, err := signing.ComputeSigningRoot(attestation, domain)
	if err != nil {
		return nil, err
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
