package inmemory

import (
	"sync"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/bloxapp/eth2-key-manager/core"
	encryptor2 "github.com/bloxapp/eth2-key-manager/encryptor"
	"github.com/bloxapp/eth2-key-manager/wallets"
)

// InMemStore implements core.Storage using in-memory store.
type InMemStore struct {
	network core.Network

	walletLock sync.Mutex
	wallet     core.Wallet

	accountsLock sync.Mutex
	accounts     map[string]*wallets.HDAccount

	highestAttestationLock sync.RWMutex
	highestAttestation     map[string]*phase0.AttestationData

	highestProposalLock sync.RWMutex
	highestProposal     map[string]*phase0.Slot

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
		accounts:           make(map[string]*wallets.HDAccount),
		highestAttestation: make(map[string]*phase0.AttestationData),
		highestProposal:    make(map[string]*phase0.Slot),
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
	store.walletLock.Lock()
	store.wallet = wallet
	store.walletLock.Unlock()
	return nil
}

// OpenWallet returns nil,nil if no wallet was found
func (store *InMemStore) OpenWallet() (core.Wallet, error) {
	if store.wallet != nil {
		store.wallet.SetContext(store.freshContext())
		return store.wallet, nil
	}
	return nil, errors.New("wallet not found")
}

// ListAccounts returns an empty array for no accounts
func (store *InMemStore) ListAccounts() ([]core.ValidatorAccount, error) {
	w, err := store.OpenWallet()
	if err != nil {
		return nil, err
	}

	return w.Accounts(), nil
}

// SaveAccount saves the given account
func (store *InMemStore) SaveAccount(account core.ValidatorAccount) error {
	store.accountsLock.Lock()
	store.accounts[account.ID().String()] = account.(*wallets.HDAccount)
	store.accountsLock.Unlock()
	return nil
}

// DeleteAccount deletes account by its ID
func (store *InMemStore) DeleteAccount(accountID uuid.UUID) error {
	store.accountsLock.Lock()
	defer store.accountsLock.Unlock()

	_, exists := store.accounts[accountID.String()]
	if !exists {
		return errors.New("account not found")
	}
	delete(store.accounts, accountID.String())
	return nil
}

// OpenAccount returns nil,nil if no account was found
func (store *InMemStore) OpenAccount(accountID uuid.UUID) (core.ValidatorAccount, error) {
	store.accountsLock.Lock()
	val := store.accounts[accountID.String()]
	store.accountsLock.Unlock()
	return val, nil
}

// SetEncryptor is the encryptor setter
func (store *InMemStore) SetEncryptor(encryptor encryptor2.Encryptor, password []byte) {
	store.encryptor = encryptor
	store.encryptionPassword = password
}

func (store *InMemStore) freshContext() *core.WalletContext {
	return &core.WalletContext{
		Storage: store,
	}
}
