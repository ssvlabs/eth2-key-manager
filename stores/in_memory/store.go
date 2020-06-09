package in_memory

import (
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
func (store *InMemStore) OpenAccount(uuid uuid.UUID) (core.Account,error) {
	if val := store.memory[uuid.String()]; val != nil {
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
	// TODO encrypt
	if val := store.memory["portfolio_seed"]; val != nil {
		return val.([]byte),nil
	} else {
		return nil,nil
	}
}

func (store *InMemStore) SecurelySavePortfolioSeed(secret []byte) error {
	// TODO decrypt
	store.memory["portfolio_seed"] = secret
	return nil
}

//func (store *InMemStore) maybeEncrypt(input interface{}) (map[string][]byte,error) {
//	data,err := json.Marshal(input)
//	if err != nil {
//		return nil,err
//	}
//
//	if store.verifyCanEncrypt() {
//		encrypted,err := input.(core.KeyBarer).EncryptedPrivateKey(store.encryptor,store.encryptionPassword)
//		if err != nil {
//			return nil,err
//		}
//		data,err = json.Marshal(encrypted)
//		if err != nil {
//			return nil,err
//		}
//	}
//
//	return data,nil
//}
//
//func (store *InMemStore) maybeDecrypt(input []byte, ret interface{}) error {
//	if store.verifyCanEncrypt() {
//		// get encrypted data
//		var data map[string]interface{}
//		err := json.Unmarshal(input,&data)
//		if err != nil {
//			return err
//		}
//
//		// decrypt
//		decrypted,err := store.encryptor.Decrypt(data,store.encryptionPassword)
//		if err != nil {
//			return err
//		}
//
//		// unmarshal to object
//		return json.Unmarshal(decrypted,&ret)
//	} else {
//		// if not encrypted just unmarshal
//		return json.Unmarshal(input,&ret)
//	}
//}

func (store *InMemStore) freshContext() *core.PortfolioContext {
	return &core.PortfolioContext {
		Storage:     store,
	}
}

func (store *InMemStore) verifyCanEncrypt() bool {
	if store.encryptor != nil {
		if store.encryptionPassword == nil {
			return false
		}
		return true
	}
	return false
}