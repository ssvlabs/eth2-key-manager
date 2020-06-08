package in_memory

import (
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func TestNonExistingPortfolio(t *testing.T) {
	stores.TestingNonExistingPortfolio(getStorage(),t)
}

func TestPortfolioStorage(t *testing.T) {
	stores.TestingPortfolioStorage(getStorage(),t)
}
