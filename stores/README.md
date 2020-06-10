# Blox KeyVault - Stores


[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)

Store is a place that saves and fetches data used by a portfolio,wallet and accounts. 
It also stores sensitive information like a seed.

Currently there are the following implementations:
- In memory storage (mostly used for testing as a quick storage setup)
- (Hashicorp's Vault)[https://www.vaultproject.io]


#### Develop you own store
You could develop you own store, for example saving it to an S3, local file system and so on.
To implement a store, simple override the methods below from [here](https://github.com/bloxapp/KeyVault/blob/master/core/storage.go)
```go
// Implements methods to store and retrieve data
// Any encryption is done on the implementation level but is not obligatory
type Storage interface {
	Name() string
	/*
		Portfolio specific
	 */
	SavePortfolio(portfolio Portfolio) error
	// will return nil,nil if no portfolio was found
	OpenPortfolio() (Portfolio,error)
	ListWallets() ([]Wallet,error)

	///*
	//	Wallet specific
	// */
	SaveWallet(wallet Wallet) error
	// will return nil,nil if no wallet was found
	OpenWallet(uuid uuid.UUID) (Wallet,error)
	// will return an empty array for no accounts
	ListAccounts(walletID uuid.UUID) ([]Account,error)

	///*
	//	Account specific
	// */
	SaveAccount(account Account) error
	// will return nil,nil if no account was found
	OpenAccount(walletId uuid.UUID, accountId uuid.UUID) (Account,error)

	// could also bee set to nil
	SetEncryptor(encryptor types.Encryptor, password []byte)
	//
	SecurelyFetchPortfolioSeed() ([]byte,error)
	//
	SecurelySavePortfolioSeed(secret []byte) error
}

```
