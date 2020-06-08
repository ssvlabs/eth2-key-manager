package core

import "github.com/google/uuid"

type WalletType = string
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
	// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
	// This will error if an account with the name already exists.
	CreateValidatorAccount(name string) (Account, error)
	// GetWithdrawalAccount returns this wallet's withdrawal key pair in the wallet as described in EIP-2334.
	// This will error if an account with the name already exists.
	GetWithdrawalAccount() (Account, error)
	// Accounts provides all accounts in the wallet.
	Accounts() <-chan Account
	// AccountByID provides a single account from the wallet given its ID.
	// This will error if the account is not found.
	// should return account = nil if not found (not an error!)
	AccountByID(id uuid.UUID) (Account, error)
	// AccountByName provides a single account from the wallet given its name.
	// This will error if the account is not found.
	// should return account = nil if not found (not an error!)
	AccountByName(name string) (Account, error)
	//
	SetContext(ctx *PortfolioContext)
}
