package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/bloxapp/eth2-key-manager/core"
)

// SignBeaconAttestation signs beacon attestation data
func (signer *SimpleSigner) SignBeaconAttestation(attestation *phase0.AttestationData, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}
	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "attestation")
	defer func() {
		signer.unlock(account.ID(), "attestation")
	}()

	// 3. far future check
	if !IsValidFarFutureEpoch(signer.network, attestation.Target.Epoch) {
		return nil, nil, errors.Errorf("target epoch too far into the future")
	}
	if !IsValidFarFutureEpoch(signer.network, attestation.Source.Epoch) {
		return nil, nil, errors.Errorf("source epoch too far into the future")
	}

	// 4. check we can even sign this
	if val, err := signer.slashingProtector.IsSlashableAttestation(pubKey, attestation); err != nil || val != nil {
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, errors.Errorf("slashable attestation (%s), not signing", val.Status)
	}

	// 5. add to protection storage
	if err := signer.slashingProtector.UpdateHighestAttestation(pubKey, attestation); err != nil {
		return nil, nil, err
	}

	// 6. Prepare and sign data
	root, err := core.ComputeETHSigningRoot(attestation, domain)
	if err != nil {
		return nil, nil, err
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}
