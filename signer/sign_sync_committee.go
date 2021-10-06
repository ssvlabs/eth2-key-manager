package signer

import (
	"encoding/hex"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
)

// SignSyncCommittee sign sync committee
func (signer *SimpleSigner) SignSyncCommittee(msgBlockRoot []byte, domain []byte, pubKey []byte) ([]byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "sync_committee")
	defer signer.unlock(account.ID(), "sync_committee")

	// 3. sign
	sszRoot := types.SSZBytes(msgBlockRoot)
	root, err := helpers.ComputeSigningRoot(&sszRoot, domain)
	if err != nil {
		return nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// SignSyncCommitteeSelectionData sign sync committee slection data
func (signer *SimpleSigner) SignSyncCommitteeSelectionData(data *eth.SyncAggregatorSelectionData, domain []byte, pubKey []byte) ([]byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "sync_committee_selection_data")
	defer signer.unlock(account.ID(), "sync_committee_selection_data")

	// 3. sign
	if data == nil {
		return nil, errors.New("selection data nil")
	}
	root, err := helpers.ComputeSigningRoot(data, domain)
	if err != nil {
		return nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// SignSyncCommitteeContributionAndProof sign sync committee
func (signer *SimpleSigner) SignSyncCommitteeContributionAndProof(contribAndProof *eth.ContributionAndProof, domain []byte, pubKey []byte) ([]byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "sync_committee_selection_and_proof")
	defer signer.unlock(account.ID(), "sync_committee_selection_and_proof")

	// 3. sign
	if contribAndProof == nil {
		return nil, errors.New("contrib proof data nil")
	}
	root, err := helpers.ComputeSigningRoot(contribAndProof, domain)
	if err != nil {
		return nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
