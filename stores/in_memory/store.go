package in_memory

import (
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	uuid "github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"reflect"
)

type InMemStore struct {
	memory         		map[string][]byte
	attMemory      		map[string]*core.BeaconAttestation
	proposalMemory 		map[string]*core.BeaconBlockHeader
	encryptor	   		types.Encryptor
	encryptionPassword 	[]byte

	portfolioType reflect.Type
	walletType reflect.Type
	accountType reflect.Type
}

func NewInMemStore(
	portfolioType reflect.Type,
	walletType reflect.Type,
	accountType reflect.Type,
	) *InMemStore {
	return NewInMemStoreWithEncryptor(nil,nil, portfolioType,walletType,accountType)
}

func NewInMemStoreWithEncryptor(
	encryptor types.Encryptor,
	password []byte,
	portfolioType reflect.Type,
	walletType reflect.Type,
	accountType reflect.Type,
	) *InMemStore {
	return &InMemStore{
		memory:         	make(map[string][]byte),
		attMemory:      	make(map[string]*core.BeaconAttestation),
		proposalMemory: 	make(map[string]*core.BeaconBlockHeader),
		encryptor:			encryptor,
		encryptionPassword:	password,
		portfolioType:portfolioType,
		walletType:walletType,
		accountType:accountType,
	}
}

// Name provides the name of the store
func (store *InMemStore) Name() string {
	return "in-memory"
}

func (store *InMemStore) SavePortfolio(portfolio core.Portfolio) error {
	data,err := store.maybeEncrypt(portfolio)
	if err != nil {
		return err
	}
	store.memory["portfolio"] = data
	return nil
}

// will return nil,nil if no portfolio was found
func (store *InMemStore) OpenPortfolio() (core.Portfolio,error) {
	bytes  := store.memory["portfolio"]
	if bytes != nil {
		ret := reflect.New(store.portfolioType.Elem()).Interface()
		err := store.maybeDecrypt(bytes,ret)
		if err != nil {
			return nil,err
		}

		// set this storage as context so p.Wallets() will work
		// TODO - this needs to be better, we shouldn't fetch portoflio to get wallets. Storage should literally be a db
		ctx := &core.PortfolioContext {
			Storage:     store,
			PortfolioId: ret.(core.Portfolio).ID(),
		}
		ret.(core.Portfolio).SetContext(ctx)

		return ret.(core.Portfolio),nil
	}
	return nil,nil
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
	data,err := store.maybeEncrypt(wallet)
	if err != nil {
		return err
	}
	store.memory[wallet.ID().String()] = data
	return nil
}

// will return nil,nil if no wallet was found
func (store *InMemStore) OpenWallet(uuid uuid.UUID) (core.Wallet,error) {
	bytes  := store.memory[uuid.String()]
	if bytes != nil {
		ret := reflect.New(store.walletType.Elem()).Interface()
		err := store.maybeDecrypt(bytes,ret)
		if err != nil {
			return nil,err
		}

		return ret.(core.Wallet),nil
	}
	return nil,nil
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
	data,err := store.maybeEncrypt(account)
	if err != nil {
		return err
	}
	store.memory[account.ID().String()] = data
	return nil
}

// will return nil,nil if no account was found
func (store *InMemStore) OpenAccount(uuid uuid.UUID) (core.Account,error) {
	bytes  := store.memory[uuid.String()]
	if bytes != nil {
		ret := reflect.New(store.accountType.Elem()).Interface()
		err := store.maybeDecrypt(bytes,ret)
		if err != nil {
			return nil,err
		}

		return ret.(core.Account),nil
	}
	return nil,nil
}

func (store *InMemStore) SetEncryptor(encryptor types.Encryptor, password []byte) {
	store.encryptor = encryptor
	store.encryptionPassword = password
}

func (store *InMemStore) maybeEncrypt(input interface{}) ([]byte,error) {
	data,err := json.Marshal(input)
	if err != nil {
		return nil,err
	}

	if store.verifyCanEncrypt() {
		encrypted,err := store.encryptor.Encrypt(data,store.encryptionPassword)
		if err != nil {
			return nil,err
		}
		data,err = json.Marshal(encrypted)
		if err != nil {
			return nil,err
		}
	}

	return data,nil
}

func (store *InMemStore) maybeDecrypt(input []byte, ret interface{}) error {
	if store.verifyCanEncrypt() {
		// get encrypted data
		var data map[string]interface{}
		err := json.Unmarshal(input,&data)
		if err != nil {
			return err
		}

		// decrypt
		decrypted,err := store.encryptor.Decrypt(data,store.encryptionPassword)
		if err != nil {
			return err
		}

		// unmarshal to object
		return json.Unmarshal(decrypted,&ret)
	} else {
		// if not encrypted just unmarshal
		return json.Unmarshal(input,&ret)
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