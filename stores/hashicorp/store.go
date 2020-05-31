package hashicorp

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

type HashicorpVaultStore struct {
	storage logical.Storage
	ctx context.Context
}

func NewHashicorpVaultStore(storage logical.Storage, ctx context.Context) *HashicorpVaultStore {
	return &HashicorpVaultStore{
		storage: storage,
		ctx:     ctx,
	}
}

const (
	WalletBasePathStr = "wallets/"
	WalletIdsBaseePath = WalletBasePathStr + "ids/"
	WalletDataPathStr = WalletIdsBaseePath + "%s/data"

	WalletAccountBase = WalletBasePathStr + "ids/%s/accounts/"
	WalletAccountPath = WalletAccountBase + "%s"


	WalletsIdMappingPathStr = WalletBasePathStr + "mappings/%s"
)

// Name provides the name of the store
func (store *HashicorpVaultStore) Name() string {
	return "Hashicorp Vault"
}


// StoreWallet stores wallet data.  It will fail if it cannot store the data.
func (store *HashicorpVaultStore) StoreWallet(walletID uuid.UUID, walletName string, data []byte) error {
	if len(walletName) == 0 {
		return fmt.Errorf("wallet name must be provided")
	}

	// put wallet data
	path := fmt.Sprintf(WalletDataPathStr, walletID.String())
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	err := store.storage.Put(store.ctx, entry)
	if err != nil {
		return err
	}

	// add map from wallet id to name
	path = fmt.Sprintf(WalletsIdMappingPathStr, walletName)
	entry = &logical.StorageEntry{
		Key:      path,
		Value:    []byte(walletID.String()),
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// RetrieveWallet retrieves wallet data for all wallets.
func (store *HashicorpVaultStore) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)

	go func() {
		walletNames,err := store.storage.List(store.ctx, WalletIdsBaseePath)
		if err == nil {
			for _, w := range walletNames {
				path := fmt.Sprintf(WalletDataPathStr, cleanFromPathSymbols(w))
				entry,error := store.storage.Get(store.ctx,path)
				if error != nil || entry == nil {
					continue
				}
				ch <- entry.Value
			}
		}
		close(ch)
	}()

	return ch
}

// RetrieveWallet retrieves wallet data for a wallet with a given name.
// It will fail if it cannot retrieve the data.
func (store *HashicorpVaultStore) RetrieveWallet(walletName string) ([]byte, error) {
	// first find the mapping between wallet name and id
	mappingpath := fmt.Sprintf(WalletsIdMappingPathStr, walletName)
	entry,error := store.storage.Get(store.ctx, mappingpath)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("could not retrieve mappings to wallet id from wallet name: %s", walletName)
	}

	// second return wallet by id
	walletId, err := uuid.Parse(string(entry.Value))
	if err != nil {
		return nil, err
	}
	return store.RetrieveWalletByID(walletId)
}

// RetrieveWalletByID retrieves wallet data for a wallet with a given ID.
// It will fail if it cannot retrieve the data.
func (store *HashicorpVaultStore) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
	path := fmt.Sprintf(WalletDataPathStr, walletID.String())
	entry,error := store.storage.Get(store.ctx,path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("wallet not found") // important as github.com/wealdtech/go-eth2-wallet-hd looks for this error
	}

	return entry.Value,nil
}

// StoreAccount stores account data.  It will fail if it cannot store the data.
// will fail for non existing wallet
func (store *HashicorpVaultStore) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
	_, err := store.RetrieveWalletByID(walletID)
	if err != nil {
		return err
	}

	// store account
	path := fmt.Sprintf(WalletAccountPath, walletID.String(), accountID.String())
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// RetrieveAccounts retrieves account information for all accounts.
func (store *HashicorpVaultStore) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
	ch := make(chan []byte, 1024)

	go func() {
		path := fmt.Sprintf(WalletAccountBase, walletID.String())
		accountNames,err := store.storage.List(store.ctx, path)
		if err == nil {
			for _, a := range accountNames {
				path := fmt.Sprintf(WalletAccountPath, walletID.String(), a)
				entry,error := store.storage.Get(store.ctx,path)
				if error != nil || entry == nil {
					continue
				}
				ch <- entry.Value
			}
		}
		close(ch)
	}()

	return ch
}

// RetrieveAccount retrieves account data for a wallet with a given ID.
// It will fail if it cannot retrieve the data.
func (store *HashicorpVaultStore) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
	path := fmt.Sprintf(WalletAccountPath, walletID.String(), accountID.String())
	entry,error := store.storage.Get(store.ctx,path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("could not retrieve account id %s for wallet id: %s", accountID.String(), walletID.String())
	}

	return entry.Value,nil
}

// StoreAccountsIndex stores the index of accounts for a given wallet.
func (store *HashicorpVaultStore) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
	return fmt.Errorf("StoreAccountsIndex not implemented")
}

// RetrieveAccountsIndex retrieves the index of accounts for a given wallet.
func (store *HashicorpVaultStore) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
	return nil,fmt.Errorf("RetrieveAccountsIndex not implemented")
}

func cleanFromPathSymbols(str string) string {
	return strings.Replace(str,"/","",-1)
}