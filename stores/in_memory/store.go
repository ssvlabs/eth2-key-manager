package in_memory

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	uuid "github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type InMemStore struct {
	memory         		map[string]interface{}
	attMemory      		map[string]*core.BeaconAttestation
	proposalMemory 		map[string]*core.BeaconBlockHeader
	encryptor	   		types.Encryptor
	encryptionPassword 	[]byte
}

func NewInMemStore(
	) *InMemStore {
	return NewInMemStoreWithEncryptor(nil,nil)
}

func NewInMemStoreWithEncryptor(
	encryptor types.Encryptor,
	password []byte,
	) *InMemStore {
	return &InMemStore{
		memory:         	make(map[string]interface{}),
		attMemory:      	make(map[string]*core.BeaconAttestation),
		proposalMemory: 	make(map[string]*core.BeaconBlockHeader),
		encryptor:			encryptor,
		encryptionPassword:	password,
	}
}

// Name provides the name of the store
func (store *InMemStore) Name() string {
	return "in-memory"
}

func (store *InMemStore) SavePortfolio(portfolio core.Portfolio) error {
	store.memory["portfolio"] = portfolio
	return nil
}

// will return nil,nil if no portfolio was found
func (store *InMemStore) OpenPortfolio() (core.Portfolio,error) {
	if val := store.memory["portfolio"]; val != nil {
		return val.(core.Portfolio),nil
	} else {
		return nil,nil
	}
}

func (store *InMemStore) ListWallets() ([]core.Wallet,error) {
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

func (store *InMemStore) SaveWallet(wallet core.Wallet) error {
	store.memory[wallet.ID().String()] = wallet
	return nil
}

// will return nil,nil if no wallet was found
func (store *InMemStore) OpenWallet(uuid uuid.UUID) (core.Wallet,error) {
	if val := store.memory[uuid.String()]; val != nil {
		return val.(core.Wallet),nil
	} else {
		return nil,nil
	}
}

// will return an empty array for no accounts
func (store *InMemStore) ListAccounts(walletID uuid.UUID) ([]core.Account,error) {
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

func (store *InMemStore) SaveAccount(account core.Account) error {
	store.memory[account.ID().String()] = account
	return nil
}

// will return nil,nil if no account was found
func (store *InMemStore) OpenAccount(walletId uuid.UUID, accountId uuid.UUID) (core.Account,error) {
	if val := store.memory[accountId.String()]; val != nil {
		return val.(core.Account),nil
	} else {
		return nil,nil
	}
}

func (store *InMemStore) SetEncryptor(encryptor types.Encryptor, password []byte) {
	store.encryptor = encryptor
	store.encryptionPassword = password
}

func (store *InMemStore) SecurelyFetchPortfolioSeed() ([]byte,error) {
	if val := store.memory["portfolio_seed"]; val != nil {
		if store.canEncrypt() {
			if encrypted,ok := val.(map[string]interface{}); ok {
				decrypted,err := store.encryptor.Decrypt(encrypted,store.encryptionPassword)
				if err != nil {
					return nil,err
				}
				return decrypted,nil
			} else {
				return nil,fmt.Errorf("no encrypted data exists")
			}

		} else {
			return val.([]byte),nil
		}
	} else {
		return nil,nil
	}
}

func (store *InMemStore) SecurelySavePortfolioSeed(secret []byte) error {
	if len(secret) != 32 {
		return fmt.Errorf("secret can be only 32 bytes (not %d bytes)",len(secret))
	}
	if store.canEncrypt() {
		encrypted,err := store.encryptor.Encrypt(secret,store.encryptionPassword)
		if err != nil {
			return err
		}
		store.memory["portfolio_seed"] = encrypted
	} else {
		store.memory["portfolio_seed"] = secret
	}
	return nil
}

func (store *InMemStore) freshContext() *core.PortfolioContext {
	return &core.PortfolioContext {
		Storage:     store,
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