package in_memory

import (
	"github.com/bloxapp/eth2-key-manager/core"
	encryptor2 "github.com/bloxapp/eth2-key-manager/encryptor"
	"github.com/bloxapp/eth2-key-manager/wallet_hd"
	uuid "github.com/google/uuid"
	"github.com/pkg/errors"
)

// InMemStore implements core.Storage using in-memory store.
type InMemStore struct {
	network            core.Network
	wallet             *wallet_hd.HDWallet
	accounts           map[string]*wallet_hd.HDAccount
	highestAttestation map[string]*core.BeaconAttestation
	proposalMemory     map[string]*core.BeaconBlockHeader
	encryptor          encryptor2.Encryptor
	encryptionPassword []byte
}

// NewInMemStore is the constructor of InMemStore.
func NewInMemStore(network core.Network) *InMemStore {
	return NewInMemStoreWithEncryptor(network, nil, nil)
}

// NewInMemStoreWithEncryptor is the constructor of InMemStore.
func NewInMemStoreWithEncryptor(network core.Network, encryptor encryptor2.Encryptor, password []byte) *InMemStore {
	return &InMemStore{
		network:            network,
		accounts:           make(map[string]*wallet_hd.HDAccount),
		highestAttestation: make(map[string]*core.BeaconAttestation),
		proposalMemory:     make(map[string]*core.BeaconBlockHeader),
		encryptor:          encryptor,
		encryptionPassword: password,
	}
}

// Name provides the name of the store.
func (store *InMemStore) Name() string {
	return "in-memory"
}

// Network returns the network.
func (store *InMemStore) Network() core.Network {
	return store.network
}

// SaveWallet implements core.Storage interface.
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
	return nil, errors.New("wallet not found")
}

// will return an empty array for no accounts
func (store *InMemStore) ListAccounts() ([]core.ValidatorAccount, error) {
	w, err := store.OpenWallet()
	if err != nil {
		return nil, err
	}

	return w.Accounts(), nil
}

func (store *InMemStore) SaveAccount(account core.ValidatorAccount) error {
	store.accounts[account.ID().String()] = account.(*wallet_hd.HDAccount)
	return nil
}

func (store *InMemStore) DeleteAccount(accountId uuid.UUID) error {
	_, exists := store.accounts[accountId.String()]
	if !exists {
		return errors.New("account not found")
	}
	delete(store.accounts, accountId.String())
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

func (store *InMemStore) SetEncryptor(encryptor encryptor2.Encryptor, password []byte) {
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
