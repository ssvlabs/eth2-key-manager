package hashicorp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bloxapp/KeyVault/stores/in_memory"

	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type HashicorpVaultStore struct {
	storage logical.Storage
	ctx     context.Context

	encryptor          types.Encryptor
	encryptionPassword []byte
}

func NewHashicorpVaultStore(storage logical.Storage, ctx context.Context) *HashicorpVaultStore {
	return &HashicorpVaultStore{
		storage: storage,
		ctx:     ctx,
	}
}

const (
	WalletDataPath = "wallet/data"

	AccountBase = "wallet/accounts/"
	AccountPath = AccountBase + "%s"
)

func FromInMemoryStore(inMem *in_memory.InMemStore, storage logical.Storage, ctx context.Context) (*HashicorpVaultStore,error) {
	// first delete old data
	// delete all accounts
	res,err := storage.List(ctx,AccountBase)
	if err != nil {
		return nil, err
	}
	for _,accountId := range res {
		path := fmt.Sprintf(AccountPath, accountId)
		err = storage.Delete(ctx, path)
		if err != nil {
			return nil, err
		}
	}
	err = storage.Delete(ctx,WalletDataPath)
	if err != nil {
		return nil, err
	}
	err = storage.Delete(ctx,AccountBase)
	if err != nil {
		return nil, err
	}

	// get new store
	newStore := NewHashicorpVaultStore(storage, ctx)

	// save wallet
	wallet,err := inMem.OpenWallet()
	if err != nil {
		return nil, err
	}
	err = newStore.SaveWallet(wallet)
	if err != nil {
		return nil, err
	}

	// save accounts
	for acc := range wallet.Accounts() {
		err = newStore.SaveAccount(acc)
		if err != nil {
			return nil, err
		}
	}

	return newStore, nil
}

func (store *HashicorpVaultStore) Name() string {
	return "Hashicorp Vault"
}

func (store *HashicorpVaultStore) SaveWallet(wallet core.Wallet) error {
	// data
	data, err := json.Marshal(wallet)
	if err != nil {
		return errors.Wrap(err, "failed to marshal wallet")
	}

	// put wallet data
	path := WalletDataPath
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}

	return store.storage.Put(store.ctx, entry)
}

// will return nil,nil if no wallet was found
func (store *HashicorpVaultStore) OpenWallet() (core.Wallet, error) {
	path := WalletDataPath
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	// un-marshal
	ret := &wallet_hd.HDWallet{} // not hardcode HDWallet
	ret.SetContext(store.freshContext())
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal HD Wallet object")
	}

	return ret, nil
}

// will return an empty array for no accounts
func (store *HashicorpVaultStore) ListAccounts() ([]core.ValidatorAccount, error) {
	w, err := store.OpenWallet()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get wallet")
	}

	ret := make([]core.ValidatorAccount, 0)
	for a := range w.Accounts() {
		ret = append(ret, a)
	}

	return ret, nil
}

func (store *HashicorpVaultStore) SaveAccount(account core.ValidatorAccount) error {
	// data
	data, err := json.Marshal(account)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account object")
	}

	// put wallet data
	path := fmt.Sprintf(AccountPath, account.ID().String())
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// will return nil,nil if no account was found
func (store *HashicorpVaultStore) OpenAccount(accountId uuid.UUID) (core.ValidatorAccount, error) {
	path := fmt.Sprintf(AccountPath, accountId)
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	// un-marshal
	ret := &wallet_hd.HDAccount{} // not hardcode HDAccount
	ret.SetContext(store.freshContext())
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal HD account object")
	}
	return ret, nil
}

// could also bee set to nil
func (store *HashicorpVaultStore) SetEncryptor(encryptor types.Encryptor, password []byte) {
	store.encryptor = encryptor
	store.encryptionPassword = password
}

func (store *HashicorpVaultStore) freshContext() *core.WalletContext {
	return &core.WalletContext{
		Storage: store,
	}
}

func (store *HashicorpVaultStore) canEncrypt() bool {
	if store.encryptor != nil {
		if store.encryptionPassword == nil {
			return false
		}
		return true
	}
	return false
}
