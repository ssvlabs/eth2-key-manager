package in_memory

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	uuid "github.com/google/uuid"
)

type InMemStore struct {
	memory         map[string]map[string][]byte
	accountIndx    map[string][]byte
	attMemory      map[string]*core.BeaconAttestation
	proposalMemory map[string]*core.BeaconBlockHeader
	mapNameToId    map[string]uuid.UUID
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		memory:         make(map[string]map[string][]byte),
		accountIndx:    make(map[string][]byte),
		mapNameToId:    make(map[string]uuid.UUID),
		attMemory:      make(map[string]*core.BeaconAttestation),
		proposalMemory: make(map[string]*core.BeaconBlockHeader),
	}
}

// Name provides the name of the store
func (store *InMemStore) Name() string {
	return "in-memory"
}

// StoreWallet stores wallet data.  It will fail if it cannot store the data.
func (store *InMemStore) StoreWallet(walletID uuid.UUID, walletName string, data []byte) error {
	if val := store.memory[walletID.String()]; val != nil { // existing wallet
		store.memory[walletID.String()]["wallet"] = data
	} else {
		store.memory[walletID.String()] = map[string][]byte { // new wallet
			"wallet":data,
		}
		store.mapNameToId[walletName] = walletID
	}
	return nil
}

// RetrieveWallet retrieves wallet data for all wallets.
func (store *InMemStore) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)

	go func() {
		for _, w := range store.memory {
			ch <- w["wallet"]
		}
		close(ch)
	}()

	return ch
}

// RetrieveWallet retrieves wallet data for a wallet with a given name.
// It will fail if it cannot retrieve the data.
func (store *InMemStore) RetrieveWallet(walletName string) ([]byte, error) {
	w, err := store.walletByName(walletName)
	if err != nil {
		return nil, err
	}
	return w["wallet"],nil
}

// RetrieveWalletByID retrieves wallet data for a wallet with a given ID.
// It will fail if it cannot retrieve the data.
func (store *InMemStore) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
	w, err := store.walletById(walletID)
	if err != nil {
		return nil,err
	}
	return w["wallet"],nil
}

// StoreAccount stores account data.  It will fail if it cannot store the data.
func (store *InMemStore) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	wallet,error := store.walletById(walletID)
	if error != nil {
		return error
	}

	wallet[accountID.String()] = data
	return nil
}

// RetrieveAccounts retrieves account information for all accounts.
func (store *InMemStore) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	ch := make(chan []byte, 1024)

	go func() {
		if wallet, err := store.walletById(walletID); err == nil {
			for key, a := range wallet {
				if key != "wallet" {
					ch <- a
				}
			}
		}

		close(ch)
	}()

	return ch
}

// RetrieveAccount retrieves account data for a wallet with a given ID.
// It will fail if it cannot retrieve the data.
func (store *InMemStore) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
	wallet,error := store.walletById(walletID)
	if error != nil {
		return nil,error
	}

	if val := wallet[accountID.String()]; val != nil {
		return val,nil
	}

	return nil, fmt.Errorf("account id %s in wallet id %s, not found",accountID.String(), walletID.String())
}

// StoreAccountsIndex stores the index of accounts for a given wallet.
func (store *InMemStore) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
	store.accountIndx[walletID.String()] = data
	return nil
}

// RetrieveAccountsIndex retrieves the index of accounts for a given wallet.
func (store *InMemStore) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
	return store.accountIndx[walletID.String()], nil
}

func (store *InMemStore) walletByName(walletName string) (map[string][]byte,error) {
	if walletId, ok := store.mapNameToId[walletName]; ok {
		return store.walletById(walletId)
	}
	return nil,fmt.Errorf("wallet not found") // important as github.com/wealdtech/go-eth2-wallet-hd looks for this error
}

func (store *InMemStore) walletById(walletID uuid.UUID) (map[string][]byte,error) {
	w := store.memory[walletID.String()]
	if w == nil {
		return nil, fmt.Errorf("wallet not found") // important as github.com/wealdtech/go-eth2-wallet-hd looks for this error
	}
	return w,nil
}