package in_memory

import (
	"fmt"
	uuid "github.com/google/uuid"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

type InMemStore struct {
	memory         map[string]*wallet
	accountIndx	   map[string][]byte
	attMemory      map[string]*pb.SignBeaconAttestationRequest
	proposalMemory map[string]*pb.SignBeaconProposalRequest
	mapIdToName    map[string]string
}

type wallet struct {
	name string
	id string
	wallet []byte
	accountIndex []byte
	accounts map[string][]byte
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		memory:         make(map[string]*wallet),
		accountIndx: 	make(map[string][]byte),
		mapIdToName:    make(map[string]string),
		attMemory:      make(map[string]*pb.SignBeaconAttestationRequest),
		proposalMemory: make(map[string]*pb.SignBeaconProposalRequest),
	}
}

// Name provides the name of the store
func (store *InMemStore) Name() string {
	return "in-memory"
}

// StoreWallet stores wallet data.  It will fail if it cannot store the data.
func (store *InMemStore) StoreWallet(walletID uuid.UUID, walletName string, data []byte) error {
	store.memory[walletName] = &wallet{
		name:     walletName,
		id:       walletID.String(),
		wallet:	  data,
		accounts: make(map[string][]byte),
	}
	store.mapIdToName[walletID.String()] = walletName
	return nil
}

// RetrieveWallet retrieves wallet data for all wallets.
func (store *InMemStore) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)

	go func() {
		for _, w := range store.memory {
			ch <- w.wallet
		}
		close(ch)
	}()

	return ch
}

// RetrieveWallet retrieves wallet data for a wallet with a given name.
// It will fail if it cannot retrieve the data.
func (store *InMemStore) RetrieveWallet(walletName string) ([]byte, error) {
	w, err := store.wallet(walletName)
	if err != nil {
		return nil, err
	}
	return w.wallet,nil
}

// RetrieveWalletByID retrieves wallet data for a wallet with a given ID.
// It will fail if it cannot retrieve the data.
func (store *InMemStore) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
	walletName := store.mapIdToName[walletID.String()]
	return store.RetrieveWallet(walletName)
}

// StoreAccount stores account data.  It will fail if it cannot store the data.
func (store *InMemStore) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	wallet,error := store.walletById(walletID)
	if error != nil {
		return error
	}

	wallet.accounts[accountID.String()] = data
	return nil
}

// RetrieveAccounts retrieves account information for all accounts.
func (store *InMemStore) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	ch := make(chan []byte, 1024)

	go func() {
		if wallet, err := store.walletById(walletID); err != nil {
			for _, a := range wallet.accounts {
				ch <- a
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

	for id, a := range wallet.accounts {
		if id == accountID.String() {
			return a, nil
		}
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

func (store *InMemStore) wallet(walletName string) (*wallet,error) {
	w := store.memory[walletName]
	if w == nil {
		return nil, fmt.Errorf("wallet not found") // important as github.com/wealdtech/go-eth2-wallet-hd looks for this error
	}
	return w,nil
}

func (store *InMemStore) walletById(walletID uuid.UUID) (*wallet,error) {
	if walletName, ok := store.mapIdToName[walletID.String()]; ok  {
		return store.wallet(walletName)
	}
	return nil, fmt.Errorf("wallet with id %s not found", walletID.String())
}