package in_memory

import (
	"testing"

	"github.com/bloxapp/eth2-key-manager/stores"
)

func TestStoringWithEncryption(t *testing.T) {
	stores.TestingWalletStorageWithEncryption(getStorage(), t)
}
