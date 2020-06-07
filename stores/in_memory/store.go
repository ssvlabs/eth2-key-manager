package in_memory

import (
	"encoding/json"
	"fmt"
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

		return ret.(core.Portfolio),nil
	}
	return nil,nil
}

//func (store *InMemStore) ListWallets() ([]core.Wallet,error) {
//	p,err := store.OpenPortfolio()
//	if err != nil {
//		return nil,err
//	}
//
//	ret := make([]core.Wallet,0)
//	for w := range p.Wallets() {
//		ret = append(ret,w)
//	}
//	return ret,nil
//}

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
//func (store *InMemStore) ListAccounts(walletID uuid.UUID) ([]core.Account,error) {
//	p,err := store.OpenPortfolio()
//	if err != nil {
//		return nil,err
//	}
//
//	w,err := p.WalletByID(walletID)
//	if err != nil {
//		return nil,err
//	}
//
//	ret := make([]core.Account,0)
//	for a := range w.Accounts() {
//		ret = append(ret,a)
//	}
//	return ret,nil
//}

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

	canEncrypt, err := store.verifyCanEncrypt()
	if err != nil {
		return nil,err
	}

	if canEncrypt {
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
	canEncrypt, err := store.verifyCanEncrypt()
	if err != nil {
		return err
	}
	if canEncrypt {
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
		return json.Unmarshal(decrypted,ret)
	} else {
		// if not encrypted just unmarshal
		return json.Unmarshal(input,&ret)
	}
}

func (store *InMemStore) verifyCanEncrypt() (bool,error) {
	if store.encryptor != nil {
		if store.encryptionPassword == nil {
			return false, fmt.Errorf("can't encrypt, missing password")
		}
		return true,nil
	}
	return false,nil
}

//// Name provides the name of the store
//func (store *InMemStore) Name() string {
//	return "in-memory"
//}
//
//// StoreWallet stores wallet data.  It will fail if it cannot store the data.
//func (store *InMemStore) StoreWallet(walletID uuid.UUID, walletName string, data []byte) error {
//	if val := store.memory[walletID.String()]; val != nil { // existing wallet
//		store.memory[walletID.String()]["wallet"] = data
//	} else {
//		store.memory[walletID.String()] = map[string][]byte { // new wallet
//			"wallet":data,
//		}
//		store.mapNameToId[walletName] = walletID
//	}
//	return nil
//}
//
//// RetrieveWallet retrieves wallet data for all wallets.
//func (store *InMemStore) RetrieveWallets() <-chan []byte {
//	ch := make(chan []byte, 1024)
//
//	go func() {
//		for _, w := range store.memory {
//			ch <- w["wallet"]
//		}
//		close(ch)
//	}()
//
//	return ch
//}
//
//// RetrieveWallet retrieves wallet data for a wallet with a given name.
//// It will fail if it cannot retrieve the data.
//func (store *InMemStore) RetrieveWallet(walletName string) ([]byte, error) {
//	w, err := store.walletByName(walletName)
//	if err != nil {
//		return nil, err
//	}
//	return w["wallet"],nil
//}
//
//// RetrieveWalletByID retrieves wallet data for a wallet with a given ID.
//// It will fail if it cannot retrieve the data.
//func (store *InMemStore) RetrieveWalletByID(walletID uuid.UUID) ([]byte, error) {
//	w, err := store.walletById(walletID)
//	if err != nil {
//		return nil,err
//	}
//	return w["wallet"],nil
//}
//
//// StoreAccount stores account data.  It will fail if it cannot store the data.
//func (store *InMemStore) StoreAccount(walletID uuid.UUID, accountID uuid.UUID, data []byte) error {
//	wallet,error := store.walletById(walletID)
//	if error != nil {
//		return error
//	}
//
//	wallet[accountID.String()] = data
//	return nil
//}
//
//// RetrieveAccounts retrieves account information for all accounts.
//func (store *InMemStore) RetrieveAccounts(walletID uuid.UUID) <-chan []byte {
//	ch := make(chan []byte, 1024)
//
//	go func() {
//		if wallet, err := store.walletById(walletID); err == nil {
//			for key, a := range wallet {
//				if key != "wallet" {
//					ch <- a
//				}
//			}
//		}
//
//		close(ch)
//	}()
//
//	return ch
//}
//
//// RetrieveAccount retrieves account data for a wallet with a given ID.
//// It will fail if it cannot retrieve the data.
//func (store *InMemStore) RetrieveAccount(walletID uuid.UUID, accountID uuid.UUID) ([]byte, error) {
//	wallet,error := store.walletById(walletID)
//	if error != nil {
//		return nil,error
//	}
//
//	if val := wallet[accountID.String()]; val != nil {
//		return val,nil
//	}
//
//	return nil, fmt.Errorf("account id %s in wallet id %s, not found",accountID.String(), walletID.String())
//}
//
//// StoreAccountsIndex stores the index of accounts for a given wallet.
//func (store *InMemStore) StoreAccountsIndex(walletID uuid.UUID, data []byte) error {
//	store.accountIndx[walletID.String()] = data
//	return nil
//}
//
//// RetrieveAccountsIndex retrieves the index of accounts for a given wallet.
//func (store *InMemStore) RetrieveAccountsIndex(walletID uuid.UUID) ([]byte, error) {
//	return store.accountIndx[walletID.String()], nil
//}
//
//func (store *InMemStore) walletByName(walletName string) (map[string][]byte,error) {
//	if walletId, ok := store.mapNameToId[walletName]; ok {
//		return store.walletById(walletId)
//	}
//	return nil,fmt.Errorf("wallet not found") // important as github.com/wealdtech/go-eth2-wallet-hd looks for this error
//}
//
//func (store *InMemStore) walletById(walletID uuid.UUID) (map[string][]byte,error) {
//	w := store.memory[walletID.String()]
//	if w == nil {
//		return nil, fmt.Errorf("wallet not found") // important as github.com/wealdtech/go-eth2-wallet-hd looks for this error
//	}
//	return w,nil
//}