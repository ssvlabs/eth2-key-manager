package core

import (
	"fmt"

	"github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// Network represents the network.
type Network string

// NetworkFromString returns network from the given string value
func NetworkFromString(n string) Network {
	switch n {
	case string(TestNetwork):
		return TestNetwork
	case string(LaunchTestNetwork):
		return LaunchTestNetwork
	case string(MainNetwork):
		return MainNetwork
	default:
		panic(fmt.Sprintf("undefined network %s", n))
	}
}

// ForkVersion returns the fork version of the network.
func (n Network) ForkVersion() []byte {
	switch n {
	case TestNetwork:
		return []byte{0, 0, 0, 1}
	case LaunchTestNetwork:
		return []byte{0, 0, 0, 2}
	case MainNetwork:
		return []byte{0, 0, 0, 3}
	default:
		panic(fmt.Sprintf("undefined network %s", n))
	}
}

// ForkVersion returns the fork version of the network.
func (n Network) FullPath(relativePath string) string {
	switch n {
	case TestNetwork:
		return BaseTestEIP2334Path + relativePath
	case LaunchTestNetwork:
		return BaseLaunchTestEIP2334Path + relativePath
	case MainNetwork:
		return BaseEIP2334Path + relativePath
	default:
		panic(fmt.Sprintf("undefined network %s", n))
	}
}

// Available networks.
const (
	// TestNetwork represents the test network.
	TestNetwork Network = "test"

	// LaunchTestNetwork represents Launch Test network.
	LaunchTestNetwork Network = "launchtest"

	// MainNetwork represents the main network.
	MainNetwork Network = "main"
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
	SetEncryptor(encryptor types.Encryptor, password []byte)
}
