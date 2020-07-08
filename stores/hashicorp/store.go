package hashicorp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bloxapp/KeyVault"
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
	PortfolioBasePath = "portfolio/"
	PortfolioDataPath = PortfolioBasePath + "/data"
	SeedPath          = PortfolioBasePath + "seed/"

	WalletBasePath     = PortfolioBasePath + "wallets/"
	WalletIdsBaseePath = WalletBasePath + "ids/"
	WalletDataPath     = WalletIdsBaseePath + "%s/data"

	AccountBase = WalletBasePath + "ids/%s/accounts/"
	AccountPath = AccountBase + "%s"
)

func (store *HashicorpVaultStore) Name() string {
	return "Hashicorp Vault"
}

func (store *HashicorpVaultStore) SavePortfolio(portfolio core.Portfolio) error {
	// data
	data, err := json.Marshal(portfolio)
	if err != nil {
		return errors.Wrap(err, "failed to marshal portfolio object")
	}

	// put wallet data
	entry := &logical.StorageEntry{
		Key:      PortfolioDataPath,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// will return nil,nil if no portfolio was found
func (store *HashicorpVaultStore) OpenPortfolio() (core.Portfolio, error) {
	entry, err := store.OpenPortfolioRaw()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open portfolio raw")
	}

	// un-marshal
	ret := &KeyVault.KeyVault{Context: store.freshContext()} // not hardcode KeyVault
	if err := json.Unmarshal(entry, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal key vault object")
	}

	return ret, nil
}

// used to fetch the raw bytes of a saved portfolio data
func (store *HashicorpVaultStore) OpenPortfolioRaw() ([]byte, error) {
	entry, err := store.storage.Get(store.ctx, PortfolioDataPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get record by portfolio data path")
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	return entry.Value, nil
}

func (store *HashicorpVaultStore) ListWallets() ([]core.Wallet, error) {
	p, err := store.OpenPortfolio()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open portfolio raw")
	}

	ret := make([]core.Wallet, 0)
	for w := range p.Wallets() {
		ret = append(ret, w)
	}
	return ret, nil
}

func (store *HashicorpVaultStore) SaveWallet(wallet core.Wallet) error {
	// data
	data, err := json.Marshal(wallet)
	if err != nil {
		return errors.Wrap(err, "failed to marshal wallet")
	}

	// put wallet data
	path := fmt.Sprintf(WalletDataPath, wallet.ID().String())
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}

	return store.storage.Put(store.ctx, entry)
}

// will return nil,nil if no wallet was found
func (store *HashicorpVaultStore) OpenWallet(uuid uuid.UUID) (core.Wallet, error) {
	path := fmt.Sprintf(WalletDataPath, uuid.String())
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
func (store *HashicorpVaultStore) ListAccounts(walletID uuid.UUID) ([]core.Account, error) {
	p, err := store.OpenPortfolio()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open portfolio raw")
	}

	w, err := p.WalletByID(walletID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get wallet by ID '%s'", walletID)
	}

	ret := make([]core.Account, 0)
	for a := range w.Accounts() {
		ret = append(ret, a)
	}

	return ret, nil
}

func (store *HashicorpVaultStore) SaveAccount(account core.Account) error {
	// data
	data, err := json.Marshal(account)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account object")
	}

	// put wallet data
	path := fmt.Sprintf(AccountPath, account.WalletID().String(), account.ID().String())
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// will return nil,nil if no account was found
func (store *HashicorpVaultStore) OpenAccount(walletId uuid.UUID, accountId uuid.UUID) (core.Account, error) {
	path := fmt.Sprintf(AccountPath, walletId, accountId)
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

//
func (store *HashicorpVaultStore) SecurelyFetchPortfolioSeed() ([]byte, error) {
	// get data
	entry, err := store.storage.Get(store.ctx, SeedPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get record with seed path")
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	// decrypt and return
	var data []byte
	if store.canEncrypt() {
		var input map[string]interface{}
		if err := json.Unmarshal(entry.Value, &input); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal data")
		}
		if input == nil {
			return nil, nil
		}

		if data, err = store.encryptor.Decrypt(input, store.encryptionPassword); err != nil {
			return nil, errors.Wrap(err, "failed to decrypt password")
		}
	} else {
		data = entry.Value
	}

	return data, nil
}

//
func (store *HashicorpVaultStore) SecurelySavePortfolioSeed(secret []byte) error {
	// data
	var data []byte
	if store.canEncrypt() {
		encrypted, err := store.encryptor.Encrypt(secret, store.encryptionPassword)
		if err != nil {
			return errors.Wrap(err, "failed to encrypt password")
		}

		if data, err = json.Marshal(encrypted); err != nil {
			return errors.Wrap(err, "failed to marshal encrypted data")
		}
	} else {
		data = secret
	}

	// save
	path := fmt.Sprintf(SeedPath)
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

func (store *HashicorpVaultStore) freshContext() *core.PortfolioContext {
	return &core.PortfolioContext{
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
