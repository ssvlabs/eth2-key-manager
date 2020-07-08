package hashicorp

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"

	"github.com/bloxapp/KeyVault/core"
)

const (
	WalletAttestationsBase      = "attestations/%s/"
	WalletAttestationPath       = WalletAttestationsBase + "%d"     // account/attestation
	WalletLatestAttestationPath = WalletAttestationsBase + "latest" // account/latest
	WalletProposalsPath         = "proposals/%s/%d/"                // account/proposal
)

func (store *HashicorpVaultStore) SaveAttestation(account core.Account, req *core.BeaconAttestation) error {
	path := fmt.Sprintf(WalletAttestationPath, account.ID().String(), req.Target.Epoch)
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal attestation request")
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

func (store *HashicorpVaultStore) RetrieveAttestation(account core.Account, epoch uint64) (*core.BeaconAttestation, error) {
	path := fmt.Sprintf(WalletAttestationPath, account.ID().String(), epoch)
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record from storage with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	var ret *core.BeaconAttestation
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal beacon attestation object")
	}

	return ret, nil
}

// both epochStart and epochEnd reflect saved attestations by their target epoch
func (store *HashicorpVaultStore) ListAttestations(account core.Account, epochStart uint64, epochEnd uint64) ([]*core.BeaconAttestation, error) {
	ret := make([]*core.BeaconAttestation, 0)

	for epoch := epochStart; epoch <= epochEnd; epoch++ {
		att, err := store.RetrieveAttestation(account, epoch)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve attestation with epoch %d", epoch)
		}

		if att != nil {
			ret = append(ret, att)
		}
	}

	return ret, nil
}

func (store *HashicorpVaultStore) SaveProposal(account core.Account, req *core.BeaconBlockHeader) error {
	path := fmt.Sprintf(WalletProposalsPath, account.ID().String(), req.Slot)
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal proposal request")
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

func (store *HashicorpVaultStore) RetrieveProposal(account core.Account, slot uint64) (*core.BeaconBlockHeader, error) {
	path := fmt.Sprintf(WalletProposalsPath, account.ID().String(), slot)
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	var ret *core.BeaconBlockHeader
	if err = json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal beacon block header object")
	}

	return ret, nil
}

func (store *HashicorpVaultStore) SaveLatestAttestation(account core.Account, req *core.BeaconAttestation) error {
	path := fmt.Sprintf(WalletLatestAttestationPath, account.ID().String())
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal beacon attestation object")
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

func (store *HashicorpVaultStore) RetrieveLatestAttestation(account core.Account) (*core.BeaconAttestation, error) {
	path := fmt.Sprintf(WalletLatestAttestationPath, account.ID().String())
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	var ret *core.BeaconAttestation
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal beacon attestation object")
	}

	return ret, nil
}
