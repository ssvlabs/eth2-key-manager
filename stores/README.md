# Eth Key Manager - Stores


Store is a place that saves and fetches data used by a portfolio,wallet and accounts. 
It also stores sensitive information like a seed.

Currently there are the following implementations:
- In memory storage (mostly used for testing as a quick storage setup)
- (Hashicorp's Vault)[https://www.vaultproject.io]


#### Develop you own store
You could develop you own store, for example saving it to an S3, local file system and so on.
To implement a store, simple override the methods below from [here](https://github.com/ssvlabs/eth2-key-manager/blob/master/core/storage.go)
```go
// Implements methods to store and retrieve data
// Any encryption is done on the implementation level but is not obligatory
type Storage interface {
	Name() string
	///*
	//	Wallet specific
	// */
	SaveWallet(wallet Wallet) error
	// will return nil,err if no wallet was found
	OpenWallet() (Wallet, error)
	// will return an empty array for no accounts
	ListAccounts() ([]ValidatorAccount, error)

	///*
	//	Account specific
	// */
	SaveAccount(account ValidatorAccount) error
	// will return nil,nil if no account was found
	OpenAccount(accountID uuid.UUID) (ValidatorAccount, error)

	// could also bee set to nil
	SetEncryptor(encryptor types.Encryptor, password []byte)
}

```
