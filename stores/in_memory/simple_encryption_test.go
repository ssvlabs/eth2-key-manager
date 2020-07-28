package in_memory

import (
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func TestStoringWithEncryption (t *testing.T) {
	stores.TestingWalletStorageWithEncryption(getStorage(),t)
}
