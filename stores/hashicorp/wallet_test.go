package hashicorp

import (
	"github.com/hashicorp/vault/sdk/logical"
)

func getStorage() logical.Storage {
	return &logical.InmemStorage{}
}

//func getWalletStorage() wtypes.Store {
//	return NewHashicorpVaultStore(getStorage(),context.Background())
//}
//
//func TestNonExistingWallet(t *testing.T) {
//	stores.TestingNonExistingWallet(getWalletStorage(),t)
//}
//
//func TestMultiWalletStorage(t *testing.T) {
//	stores.TestingWalletListing(getWalletStorage(),t)
//}
//
//func TestWalletStorage(t *testing.T) {
//	stores.TestingWalletStorage(getWalletStorage(),t)
//}
//
//func TestUpdatingWallet(t *testing.T) {
//	stores.TestingUpdatingWallet(getWalletStorage(),t)
//}