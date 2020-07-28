package core

import (
	"github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// Implements methods to store and retrieve data
// Any encryption is done on the implementation level but is not obligatory
type Storage interface {
	Name() string
	///*
	//	Wallet specific
	// */
	SaveWallet(wallet Wallet) error
	// will return nil,nil if no wallet was found
	OpenWallet() (Wallet,error)
	// will return an empty array for no accounts
	ListAccounts() ([]ValidatorAccount,error)

	///*
	//	Account specific
	// */
	SaveAccount(account ValidatorAccount) error
	// will return nil,nil if no account was found
	OpenAccount(accountId uuid.UUID) (ValidatorAccount,error)

	// could also bee set to nil
	SetEncryptor(encryptor types.Encryptor, password []byte)
}