package core

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/bloxapp/eth2-key-manager/encryptor"
)

var (
	// This is the testNet genesis time, 2020-11-18 12:00:07 UTC
	testNetGenesisTime = time.Date(2020, 11, 18, 12, 0, 7, 0, time.UTC)

	// This is the mainNet genesis time, 2020-12-01 12:00:00 UTC
	mainNetGenesisTime = time.Date(2020, 12, 1, 12, 0, 23, 0, time.UTC)
)

// Network represents the network.
type Network string

// NetworkFromString returns network from the given string value
func NetworkFromString(n string) Network {
	switch n {
	case string(TestNetwork):
		return TestNetwork
	case string(MainNetwork):
		return MainNetwork
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return ""
	}
}

// ForkVersion returns the fork version of the network.
func (n Network) ForkVersion() []byte {
	switch n {
	case TestNetwork:
		return []byte{0, 0, 32, 9}
	case MainNetwork:
		return []byte{0, 0, 0, 0}
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return nil
	}
}

// GenesisTime returns the genesis time of the network.
func (n Network) GenesisTime() time.Time {
	switch n {
	case TestNetwork:
		return testNetGenesisTime
	case MainNetwork:
		return mainNetGenesisTime
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return time.Time{}
	}
}

// DepositContractAddress returns the deposit contract address of the network.
func (n Network) DepositContractAddress() string {
	switch n {
	case TestNetwork:
		return "0x8c5fecdC472E27Bc447696F431E425D02dd46a8c"
	case MainNetwork:
		return "0x00000000219ab540356cBB839Cbe05303d7705Fa"
	default:
		logrus.WithField("network", n).Fatal("undefined network")
		return ""
	}
}

// ForkVersion returns the fork version of the network.
func (n Network) FullPath(relativePath string) string {
	return BaseEIP2334Path + relativePath
}

// Available networks.
const (
	// TestNetwork represents the Pyrmont test network.
	TestNetwork Network = "pyrmont"

	// MainNetwork represents the main network.
	MainNetwork Network = "mainnet"
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
