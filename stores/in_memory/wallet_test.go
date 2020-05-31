package in_memory

import (
	"github.com/bloxapp/KeyVault/stores"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
	"testing"
)

func getWalletStorage() wtypes.Store {
	return NewInMemStore()
}

func TestNonExistingWallet(t *testing.T) {
	stores.TestingNonExistingWallet(getWalletStorage(),t)
}

func TestMultiWalletStorage(t *testing.T) {
	stores.TestingMultiWalletStorage(getWalletStorage(),t)
}

func TestWalletStorage(t *testing.T) {
	stores.TestingWalletStorage(getWalletStorage(),t)
}
