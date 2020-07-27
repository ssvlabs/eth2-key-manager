package in_memory

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func getStorage() core.Storage {
	return NewInMemStore()
}

func TestOpeningAccounts(t *testing.T) {
	stores.TestingOpenAccounts(getStorage(),t)
}

func TestNonExistingWallet(t *testing.T) {
	stores.TestingNonExistingWallet(getStorage(),t)
}

func TestWalletStorage(t *testing.T) {
	stores.TestingWalletStorage(getStorage(),t)
}

