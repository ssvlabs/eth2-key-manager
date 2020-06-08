package in_memory

import (
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"reflect"
	"testing"
)

func getStorage() core.PortfolioStorage {
	return NewInMemStore(
		reflect.TypeOf(&KeyVault.KeyVault{}),
		reflect.TypeOf(&wallet_hd.HDWallet{}),
		reflect.TypeOf(&wallet_hd.HDAccount{}),
		)
}

func TestWithdrawalAccount(t *testing.T) {
	stores.TestingWithdrawalAccount(getStorage(),t)
}

func TestOpeningAccounts(t *testing.T) {
	stores.TestingOpenAccounts(getStorage(),t)
}

func TestNonExistingWallet(t *testing.T) {
	stores.TestingNonExistingWallet(getStorage(),t)
}

func TestWalletListingStorage(t *testing.T) {
	stores.TestingWalletListing(getStorage(),t)
}

func TestWalletStorage(t *testing.T) {
	stores.TestingWalletStorage(getStorage(),t)
}

