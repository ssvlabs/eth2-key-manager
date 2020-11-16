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
	case string(ZinkenNetwork):
		return ZinkenNetwork
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
	case ZinkenNetwork:
		return []byte{0, 0, 0, 3}
	case MainNetwork:
		return []byte{0, 0, 0, 4}
	default:
		panic(fmt.Sprintf("undefined network %s", n))
	}
}

// DepositContractAddress returns the deposit contract address of the network.
func (n Network) DepositContractAddress() string {
	switch n {
	case TestNetwork:
		return "0x07b39F4fDE4A38bACe212b546dAc87C58DfE3fDC"
	case ZinkenNetwork:
		return "0x99F0Ec06548b086E46Cb0019C78D0b9b9F36cD53"
	case MainNetwork:
		return "0x00000000219ab540356cBB839Cbe05303d7705Fa"
	default:
		panic(fmt.Sprintf("undefined network %s", n))
	}
}

// ForkVersion returns the fork version of the network.
func (n Network) FullPath(relativePath string) string {
	return BaseEIP2334Path + relativePath
}

// Available networks.
const (
	// TestNetwork represents the test network.
	TestNetwork Network = "test"

	// ZinkenNetwork represents Zinken network.
	ZinkenNetwork Network = "zinken"

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
	SetEncryptor(encryptor types.Encryptor, password []byte)
}
