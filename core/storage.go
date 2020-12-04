package core

import (
	"github.com/bloxapp/eth2-key-manager/encryptor"
	"github.com/google/uuid"
)

// Storage represents storage behavior
// Any encryption is done on the implementation level but is not obligatory
type Storage interface {
	WalletStorage
	AccountStorage

	// Name returns storage name.
	Name() string

	// Network returns the network storage is related to.
	Network() Network
}

// WalletStorage represents the behavior of the wallet storage
type WalletStorage interface {
	// SaveWallet stores the given wallet.
	SaveWallet(wallet Wallet) error

	// OpenWallet returns nil,err if no wallet was found
	OpenWallet() (Wallet, error)

	// ListAccounts returns an empty array for no accounts
	ListAccounts() ([]ValidatorAccount, error)
}

// AccountStorage represents the behavior of the account storage
type AccountStorage interface {
	// SaveAccount saves the given account
	SaveAccount(account ValidatorAccount) error

	// DeleteAccount deletes account by uuid
	DeleteAccount(accountID uuid.UUID) error

	// OpenAccount returns nil,nil if no account was found
	OpenAccount(accountID uuid.UUID) (ValidatorAccount, error)

	// SetEncryptor sets the given encryptor to the wallet.
	SetEncryptor(encryptor encryptor.Encryptor, password []byte)
}
