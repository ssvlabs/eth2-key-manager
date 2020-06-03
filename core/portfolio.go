package core

import (
	"github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// A portfolio is a container of wallets
type Portfolio interface {
	// CreateAccount creates a new account in the wallet.
	// This will error if an account with the name already exists.
	CreateWallet(name string) (Wallet, error)
	// Accounts provides all accounts in the wallet.
	Wallets() (<-chan Wallet,error)
	// AccountByID provides a single account from the wallet given its ID.
	// This will error if the account is not found.
	WalletByID(id uuid.UUID) (Wallet, error)
	// AccountByName provides a single account from the wallet given its name.
	// This will error if the account is not found.
	WalletByName(name string) (Wallet, error)
}

type PortfolioContext struct {
	Storage PortfolioStorage
	Encryptor types.Encryptor
	LockPassword []byte // only used internally for quick lock
}
