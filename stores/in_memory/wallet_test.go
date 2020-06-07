package in_memory

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func getStorage() core.PortfolioStorage {
	return NewInMemStore()
}

func TestNonExistingWallet(t *testing.T) {
	stores.TestingNonExistingWallet(getStorage(),t)
}

func TestMultiWalletStorage(t *testing.T) {
	stores.TestingWalletListing(getStorage(),t)
}

func TestWalletStorage(t *testing.T) {
	stores.TestingWalletStorage(getStorage(),t)
}

//func TestUpdatingWallet(t *testing.T) {
//	stores.TestingUpdatingWallet(getStorage(),t)
//}
