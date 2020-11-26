package core

import (
	"github.com/bloxapp/eth2-key-manager/encryptor"
	"github.com/google/uuid"
)

// Implements methods to store and retrieve data
// Any encryption is done on the implementation level but is not obligatory
type Storage interface {
	// Name returns storage name.
	Name() string

	// Network returns the network storage is related to.
	Network() Network

	//-------------------------
	//	Wallet specific
	//-------------------------
	// SaveWallet stores the given wallet.
	SaveWallet(wallet Wallet) error
	// OpenWallet returns nil,err if no wallet was found
	OpenWallet() (Wallet, error)
	// ListAccounts returns an empty array for no accounts
	ListAccounts() ([]ValidatorAccount, error)

	//-------------------------
	//	Account specific
	//-------------------------
	SaveAccount(account ValidatorAccount) error
	// Delete account by uuid
	DeleteAccount(accountId uuid.UUID) error
	// will return nil,nil if no account was found
	OpenAccount(accountId uuid.UUID) (ValidatorAccount, error)

	// SetEncryptor sets the given encryptor to the wallet.
	SetEncryptor(encryptor encryptor.Encryptor, password []byte)
}
