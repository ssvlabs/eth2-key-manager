package in_memory

import (
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func TestPortfolioStorage(t *testing.T) {
	stores.TestingPortfolioStorage(getStorage(),t)
}
