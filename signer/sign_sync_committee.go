package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SignSyncCommittee sign sync committee
func (signer *SimpleSigner) SignSyncCommittee(msgBlockRoot []byte, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	// 2. lock for current account
	val := signer.lock(account.ID(), "sync_committee")
	val.Lock()
	defer val.Unlock()

	// 3. sign
	sszRoot := SSZBytes(msgBlockRoot)
	root, err := ComputeETHSigningRoot(&sszRoot, domain)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}

// SignSyncCommitteeSelectionData sign sync committee slection data
func (signer *SimpleSigner) SignSyncCommitteeSelectionData(data *altair.SyncAggregatorSelectionData, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	// 2. lock for current account
	val := signer.lock(account.ID(), "sync_committee_selection_data")
	val.Lock()
	defer val.Unlock()

	// 3. sign
	if data == nil {
		return nil, nil, errors.New("selection data nil")
	}
	root, err := ComputeETHSigningRoot(data, domain)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}

// SignSyncCommitteeContributionAndProof sign sync committee
func (signer *SimpleSigner) SignSyncCommitteeContributionAndProof(contribAndProof *altair.ContributionAndProof, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	// 2. lock for current account
	val := signer.lock(account.ID(), "sync_committee_selection_and_proof")
	val.Lock()
	defer val.Unlock()

	// 3. sign
	if contribAndProof == nil {
		return nil, nil, errors.New("contrib proof data nil")
	}
	root, err := ComputeETHSigningRoot(contribAndProof, domain)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get signing root")
	}
	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}
