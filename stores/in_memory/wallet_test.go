package in_memory

import (
	"testing"

	"github.com/bloxapp/eth-key-manager/core"
	"github.com/bloxapp/eth-key-manager/stores"
)

func getStorage() core.Storage {
	return NewInMemStore()
}

func TestOpeningAccounts(t *testing.T) {
	stores.TestingOpenAccounts(getStorage(), t)
}

func TestNonExistingWallet(t *testing.T) {
	stores.TestingNonExistingWallet(getStorage(), t)
}

func TestWalletStorage(t *testing.T) {
	stores.TestingWalletStorage(getStorage(), t)
}
