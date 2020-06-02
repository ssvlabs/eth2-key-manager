package core

import "github.com/google/uuid"

type WalletType string
const (
	HDWallet WalletType = "HD" // hierarchical deterministic wallet
	ND		 WalletType = "ND" // non - deterministic
)

// A wallet is a container of accounts.
// Accounts = key pairs
type Wallet interface {
	// ID provides the ID for the wallet.
	ID() uuid.UUID
	// Name provides the name for the wallet.
	Name() string
	// Type provides the type of the wallet.
	Type() WalletType
	// CreateAccount creates a new account in the wallet.
	// This will error if an account with the name already exists.
	CreateAccount(name string) (Account, error)
	// Accounts provides all accounts in the wallet.
	Accounts() <-chan Account
	// AccountByID provides a single account from the wallet given its ID.
	// This will error if the account is not found.
	AccountByID(id uuid.UUID) (Account, error)
	// AccountByName provides a single account from the wallet given its name.
	// This will error if the account is not found.
	AccountByName(name string) (Account, error)
}
