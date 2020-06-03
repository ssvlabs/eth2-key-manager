package core

import "github.com/google/uuid"

// Implements methods to store and retrieve portfolios
// Any encryption is done on the implementation level but is not obligatory
type PortfolioStorage interface {
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
	ListAccounts() ([]*Account,error)

	///*
	//	Account specific
	// */
	SaveAccount(account Account) error
	// will return nil,nil if no account was found
	OpenAccount(uuid uuid.UUID) (Account,error)
}