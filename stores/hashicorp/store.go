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
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type HashicorpVaultStore struct {
	storage logical.Storage
	ctx context.Context

	encryptor	   		types.Encryptor
	encryptionPassword 	[]byte
}

func NewHashicorpVaultStore(storage logical.Storage, ctx context.Context) *HashicorpVaultStore {
	return &HashicorpVaultStore{
		storage: storage,
		ctx:     ctx,
	}
}

const (
	PortfolioBasePath 		= "portfolio/"
	PortfolioDataPath		= PortfolioBasePath		+ "/data"
	SeedPath				= PortfolioBasePath  	+ "seed/"

	WalletBasePath     		= PortfolioBasePath 	+ "wallets/"
	WalletIdsBaseePath 		= WalletBasePath 		+ "ids/"
	WalletDataPath     		= WalletIdsBaseePath 	+ "%s/data"

	AccountBase 			= WalletBasePath 		+ "ids/%s/accounts/"
	AccountPath 			= AccountBase 			+ "%s"
)

func (store *HashicorpVaultStore) Name() string {
	return "Hashicorp Vault"
}

func (store *HashicorpVaultStore) SavePortfolio(portfolio core.Portfolio) error {
	// data
	data,err := json.Marshal(portfolio)
	if err != nil {
		return err
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
func (store *HashicorpVaultStore) OpenPortfolio() (core.Portfolio,error) {
	entry,error := store.storage.Get(store.ctx,PortfolioDataPath)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, nil
	}

	// un-marshal
	ret := &KeyVault.KeyVault{Context:store.freshContext()} // not hardcode KeyVault
	error = json.Unmarshal(entry.Value,&ret)
	if error != nil {
		return nil,error
	}
	return ret,nil
}

func (store *HashicorpVaultStore) ListWallets() ([]core.Wallet,error) {
	p,err := store.OpenPortfolio()
	if err != nil {
		return nil,err
	}

	ret := make([]core.Wallet,0)
	for w := range p.Wallets() {
		ret = append(ret,w)
	}
	return ret,nil
}

func (store *HashicorpVaultStore) SaveWallet(wallet core.Wallet) error {
	// data
	data,err := json.Marshal(wallet)
	if err != nil {
		return err
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
func (store *HashicorpVaultStore) OpenWallet(uuid uuid.UUID) (core.Wallet,error) {
	path := fmt.Sprintf(WalletDataPath, uuid.String())
	entry,error := store.storage.Get(store.ctx,path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, nil
	}

	// un-marshal
	ret := &wallet_hd.HDWallet{} // not hardcode HDWallet
	ret.SetContext(store.freshContext())
	error = json.Unmarshal(entry.Value,&ret)
	if error != nil {
		return nil,error
	}
	return ret,nil
}

// will return an empty array for no accounts
func (store *HashicorpVaultStore) ListAccounts(walletID uuid.UUID) ([]core.Account,error) {
	p,err := store.OpenPortfolio()
	if err != nil {
		return nil,err
	}

	w,err := p.WalletByID(walletID)
	if err != nil {
		return nil,err
	}

	ret := make([]core.Account,0)
	for a := range w.Accounts() {
		ret = append(ret,a)
	}
	return ret,nil
}

func (store *HashicorpVaultStore) SaveAccount(account core.Account) error {
	// data
	data,err := json.Marshal(account)
	if err != nil {
		return err
	}

	// put wallet data
	path := fmt.Sprintf(AccountPath,account.WalletID().String(), account.ID().String())
	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// will return nil,nil if no account was found
func (store *HashicorpVaultStore) OpenAccount(walletId uuid.UUID, accountId uuid.UUID) (core.Account,error) {
	path := fmt.Sprintf(AccountPath, walletId, accountId)
	entry,error := store.storage.Get(store.ctx,path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, nil
	}

	// un-marshal
	ret := &wallet_hd.HDAccount{} // not hardcode HDAccount
	ret.SetContext(store.freshContext())
	error = json.Unmarshal(entry.Value,&ret)
	if error != nil {
		return nil,error
	}
	return ret,nil
}


// could also bee set to nil
func (store *HashicorpVaultStore) SetEncryptor(encryptor types.Encryptor, password []byte) {
	store.encryptor = encryptor
	store.encryptionPassword = password
}

//
func (store *HashicorpVaultStore) SecurelyFetchPortfolioSeed() ([]byte,error) {
	// get data
	path := fmt.Sprintf(SeedPath)
	entry,error := store.storage.Get(store.ctx,path)
	if error != nil {
		return nil, error
	}
	if entry == nil {
		return nil, nil
	}

	// decrypt and return
	var data []byte
	if store.canEncrypt() {
		var input map[string]interface{}
		err := json.Unmarshal(entry.Value,&input)
		if err != nil {
			return nil, error
		}
		if input == nil {
			return nil,nil
		}
		data,err = store.encryptor.Decrypt(input,store.encryptionPassword)
		if err != nil {
			return nil,err
		}
	} else {
		data = entry.Value
	}

	return data,nil
}

//
func (store *HashicorpVaultStore) SecurelySavePortfolioSeed(secret []byte) error {
	// data
	var data []byte
	if store.canEncrypt() {
		encrypted,err := store.encryptor.Encrypt(secret,store.encryptionPassword)
		if err != nil {
			return err
		}
		data,err = json.Marshal(encrypted)
		if err != nil {
			return err
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
	return &core.PortfolioContext {
		Storage:     store,
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