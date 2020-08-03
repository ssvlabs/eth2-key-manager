package in_memory

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/wallet_hd"
	uuid "github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type InMemStore struct {
	wallet             *wallet_hd.HDWallet
	accounts           map[string]*wallet_hd.HDAccount
	attMemory          map[string]*core.BeaconAttestation
	proposalMemory     map[string]*core.BeaconBlockHeader
	encryptor          types.Encryptor
	encryptionPassword []byte
}

func NewInMemStore() *InMemStore {
	return NewInMemStoreWithEncryptor(nil, nil)
}

func NewInMemStoreWithEncryptor(encryptor types.Encryptor, password []byte) *InMemStore {
	return &InMemStore{
		accounts:           make(map[string]*wallet_hd.HDAccount),
		attMemory:          make(map[string]*core.BeaconAttestation),
		proposalMemory:     make(map[string]*core.BeaconBlockHeader),
		encryptor:          encryptor,
		encryptionPassword: password,
	}
}

// Name provides the name of the store
func (store *InMemStore) Name() string {
	return "in-memory"
}

func (store *InMemStore) SaveWallet(wallet core.Wallet) error {
	store.wallet = wallet.(*wallet_hd.HDWallet)
	return nil
}

// will return nil,nil if no wallet was found
func (store *InMemStore) OpenWallet() (core.Wallet, error) {
	if store.wallet != nil {
		store.wallet.SetContext(store.freshContext())
		return store.wallet, nil
	}
	return nil, fmt.Errorf("wallet not found")
}

// will return an empty array for no accounts
func (store *InMemStore) ListAccounts() ([]core.ValidatorAccount, error) {
	w, err := store.OpenWallet()
	if err != nil {
		return nil, err
	}

	ret := make([]core.ValidatorAccount, 0)
	for a := range w.Accounts() {
		ret = append(ret, a)
	}
	return ret, nil
}

func (store *InMemStore) SaveAccount(account core.ValidatorAccount) error {
	store.accounts[account.ID().String()] = account.(*wallet_hd.HDAccount)
	return nil
}

// will return nil,nil if no account was found
func (store *InMemStore) OpenAccount(accountId uuid.UUID) (core.ValidatorAccount, error) {
	if val := store.accounts[accountId.String()]; val != nil {
		return val, nil
	} else {
		return nil, nil
	}
}

func (store *InMemStore) SetEncryptor(encryptor types.Encryptor, password []byte) {
	store.encryptor = encryptor
	store.encryptionPassword = password
}

func (store *InMemStore) freshContext() *core.WalletContext {
	return &core.WalletContext{
		Storage: store,
	}
}

func (store *InMemStore) canEncrypt() bool {
	if store.encryptor != nil {
		if store.encryptionPassword == nil {
			return false
		}
		return true
	}
	return false
}
