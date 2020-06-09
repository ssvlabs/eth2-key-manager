package hashicorp

import (
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func TestAccountIndexes(t *testing.T) {
	store := getWalletStorage()
	stores.TestingAccountIndexes(store,t)
}
