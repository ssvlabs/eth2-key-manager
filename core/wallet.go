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
	// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
	// This will error if an account with the name already exists.
	CreateValidatorAccount(name string) (Account, error)
	// CreateWithdrawalKey creates a new withdrawal key pair in the wallet.
	// This will error if an account with the name already exists.
	// according to EIP 2334 there is 1 withdrawal key per wallet hierarchy
	GetWithdrawalAccount() (Account, error)
	// Accounts provides all accounts in the wallet.
	Accounts() <-chan Account
	// AccountByID provides a single account from the wallet given its ID.
	// This will error if the account is not found.
	AccountByID(id uuid.UUID) (Account, error)
	// AccountByName provides a single account from the wallet given its name.
	// This will error if the account is not found.
	AccountByName(name string) (Account, error)
	// lock will encrypt the seed, save it to memory and nil the plain text seed.
	// it will use an internally save locking password so it could be locked at all times
	Lock() error
	IsLocked() bool
	// unlock will decrypt the seed and save on memory
	// it needs a provided password
	Unlock(password []byte) error
}
