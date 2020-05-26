package stores

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

type HashcorpVaultStore struct {
	storage logical.Storage
	ctx context.Context
}

const (
	WalletBasePathStr = "wallets/"
	WalletDataPathStr = WalletBasePathStr + "%s/data"

	WalletAccountBase = WalletBasePathStr + "%s/accounts/"
	WalletAccountPath = WalletAccountBase + "%s"


	WalletsIdMappingPathStr = WalletBasePathStr + "mappings/%s"
)

// StoreWallet stores wallet data.  It will fail if it cannot store the data.
func (store *HashcorpVaultStore) StoreWallet(walletID uuid.UUID, walletName string, data []byte) error {
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
func (store *HashcorpVaultStore) RetrieveWallets() <-chan []byte {
	ch := make(chan []byte, 1024)

	go func() {
		walletNames,err := store.storage.List(store.ctx, WalletBasePathStr)
		if err == nil {
			for _, w := range walletNames {
				path := fmt.Sprintf(WalletDataPathStr, w)
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
func (store *HashcorpVaultStore) RetrieveWallet(walletName string) ([]byte, error) {
	// first find the mapping between wallet name and id
	path := fmt.Sprintf(WalletsIdMappingPathStr, walletName)
	entry,error := store.storage.Get(store.ctx,path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("could not retrieve mappings to wallet id from wallet name: %s", walletName)
	}

	return entry.Value,nil
}

// RetrieveWalletByID retrieves wallet data for a wallet with a given ID.
// It will fail if it cannot retrieve the data.
func (store *HashcorpVaultStore) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
	path := fmt.Sprintf(WalletDataPathStr, walletID.String())
	entry,error := store.storage.Get(store.ctx,path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, fmt.Errorf("could not retrieve wallet id: %s", walletID.String())
	}

	return entry.Value,nil
}

// StoreAccount stores account data.  It will fail if it cannot store the data.
func (store *HashcorpVaultStore) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
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
func (store *HashcorpVaultStore) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
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
func (store *HashcorpVaultStore) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
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
func (store *HashcorpVaultStore) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
	return fmt.Errorf("StoreAccountsIndex not implemented")
}

// RetrieveAccountsIndex retrieves the index of accounts for a given wallet.
func (store *HashcorpVaultStore) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
	return nil,fmt.Errorf("RetrieveAccountsIndex not implemented")
}
