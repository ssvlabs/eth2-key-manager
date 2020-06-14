package core

import (
	"github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

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
	// used to fetch the raw bytes of a saved portfolio data
	OpenPortfolioRaw() ([]byte,error)
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