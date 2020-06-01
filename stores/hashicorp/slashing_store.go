package hashicorp


import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/hashicorp/vault/sdk/logical"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

const (
	WalletAttestationsBase = "attestations/%s/"
	WalletAttestationPath = WalletAttestationsBase + "%d" // account/attestation
	WalletLatestAttestationPath = WalletAttestationsBase + "latest" // account/latest
	WalletProposalsPath = "proposals/%s/%d/" // account/proposal
)

func (store *HashicorpVaultStore) SaveAttestation(account types.Account, req *core.BeaconAttestation) error {
	path := fmt.Sprintf(WalletAttestationPath, account.ID().String(),req.Target.Epoch)
	data,err := json.Marshal(req)
	if err != nil {
		return err
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

func (store *HashicorpVaultStore) RetrieveAttestation(account types.Account, epoch uint64) (*core.BeaconAttestation, error) {
	path := fmt.Sprintf(WalletAttestationPath, account.ID().String(),epoch)
	entry,error := store.storage.Get(store.ctx, path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("attestation not found")
	}

	var ret *core.BeaconAttestation
	error = json.Unmarshal(entry.Value,&ret)
	if error != nil {
		return nil, error
	}

	return ret,nil
}

// both epochStart and epochEnd reflect saved attestations by their target epoch
func (store *HashicorpVaultStore) ListAttestations(account types.Account, epochStart uint64, epochEnd uint64) ([]*core.BeaconAttestation, error) {
	ret := make([]*core.BeaconAttestation,0)
	for i:= epochStart ; i <= epochEnd ; i++ {
		att,err := store.RetrieveAttestation(account,i)
		if err != nil {
			continue
		}

		ret = append(ret,att)

	}
	return ret,nil
}

func (store *HashicorpVaultStore) SaveProposal(account types.Account, req *core.BeaconBlockHeader) error {
	path := fmt.Sprintf(WalletProposalsPath, account.ID().String(),req.Slot)
	data,err := json.Marshal(req)
	if err != nil {
		return err
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

func (store *HashicorpVaultStore) RetrieveProposal(account types.Account, slot uint64) (*core.BeaconBlockHeader, error) {
	path := fmt.Sprintf(WalletProposalsPath, account.ID().String(),slot)
	entry,error := store.storage.Get(store.ctx, path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("proposal not found")
	}

	var ret *core.BeaconBlockHeader
	error = json.Unmarshal(entry.Value,&ret)
	if error != nil {
		return nil, error
	}

	return ret,nil
}

func (store *HashicorpVaultStore) SaveLatestAttestation(account types.Account, req *core.BeaconAttestation) error {
	path := fmt.Sprintf(WalletLatestAttestationPath, account.ID().String())
	data,err := json.Marshal(req)
	if err != nil {
		return err
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

func (store *HashicorpVaultStore) RetrieveLatestAttestation(account types.Account) (*core.BeaconAttestation, error) {
	path := fmt.Sprintf(WalletLatestAttestationPath, account.ID().String())
	entry,error := store.storage.Get(store.ctx, path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("attestation not found")
	}

	var ret *core.BeaconAttestation
	error = json.Unmarshal(entry.Value,&ret)
	if error != nil {
		return nil, error
	}

	return ret,nil
}
