package core

// Implements methods to store and retrieve portfolios
// Any encryption is done on the implementation level but is not obligatory
type PortfolioStorage interface {
	/*
		Portfolio specific
	 */
	SavePortfolio(portfolio Portfolio) error
	// will return nil,nil if no portfolio was found
	OpenPortfolio() (Portfolio,error)

	///*
	//	Wallet specific
	// */
	//StoreWallet(wallet Wallet) error
	//// will return nil,nil if no wallet was found
	//GetWallet(uuid uuid.UUID) (Wallet,error)
	//// will return nil,nil if no wallet was found
	//GetWalletByName(name string) (Wallet,error)
	//// will return an empty array for no wallets
	//ListWallets() ([]*Wallet,error)
	//
	///*
	//	Account specific
	// */
	//StoreAccount(account Account) error
	//// will return nil,nil if no account was found
	//GetAccount(uuid uuid.UUID) (Account,error)
	//// will return nil,nil if no account was found
	//GetAccountByName(name string) (Account,error)
	//// will return an empty array for no accounts
	//ListAccounts() ([]*Account,error)
}