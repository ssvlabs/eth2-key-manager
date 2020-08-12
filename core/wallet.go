package core

import "github.com/google/uuid"

type WalletType = string

const (
	HDWallet WalletType = "HD" // hierarchical deterministic wallet
	ND       WalletType = "ND" // non - deterministic
)

// A wallet is a container of accounts.
// Accounts = key pairs
type Wallet interface {
	// ID provides the ID for the wallet.
	ID() uuid.UUID
	// Type provides the type of the wallet.
	Type() WalletType
	// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
	// This will error if an account with the name already exists.
	CreateValidatorAccount(seed []byte, name string) (ValidatorAccount, error)
	// Accounts provides all accounts in the wallet.
	Accounts() <-chan ValidatorAccount
	// AccountByID provides a single account from the wallet given its ID.
	// This will error if the account is not found.
	// should return account = nil if not found (not an error!)
	AccountByID(id uuid.UUID) (ValidatorAccount, error)
	// AccountByPublicKey provides a single account from the wallet given its public key.
	// This will error if the account is not found.
	// should return account = nil if not found (not an error!)
	AccountByPublicKey(pubKey string) (ValidatorAccount, error)
	// DeleteAccountByPublicKey delete an account from the wallet given its public key.
	// This will error if the account is not found.
	// should return nil if not error otherwise the error
	DeleteAccountByPublicKey(pubKey string) error
	SetContext(ctx *WalletContext)
}

type WalletContext struct {
	Storage Storage
}
