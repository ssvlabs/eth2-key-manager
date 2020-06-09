package hashicorp

//func getPopulatedWalletStorage(t *testing.T) (wtypes.Store,[]uuid.UUID) {
//	store := getWalletStorage()
//
//	ids := []uuid.UUID{
//		uuid.New(),
//		uuid.New(),
//		uuid.New(),
//		uuid.New(),
//	}
//
//	if err := store.StoreWallet(ids[0],"1", []byte("wallet 1")); err != nil { t.Error(err) }
//	if err := store.StoreWallet(ids[1],"2", []byte("wallet 2")); err != nil { t.Error(err) }
//	if err := store.StoreWallet(ids[2],"3", []byte("wallet 3")); err != nil { t.Error(err) }
//	if err := store.StoreWallet(ids[3],"4", []byte("wallet 4")); err != nil { t.Error(err) }
//
//	return store,ids
//}
//
//func TestAccountAtNonExistingWallet(t *testing.T) {
//	storage, ids := getPopulatedWalletStorage(t)
//	stores.TestingAccountAtNonExistingWallet(storage, ids, t)
//}
//
//func TestAddingAccountsToWallet(t *testing.T) {
//	storage, ids := getPopulatedWalletStorage(t)
//	stores.TestingAddingAccountsToWallet(storage, ids, t)
//}
//
//func TestFetchingNonExistingAccount(t *testing.T) {
//	storage, ids := getPopulatedWalletStorage(t)
//	stores.TestingFetchingNonExistingAccount(storage, ids, t)
//}
//
//func TestListingAccounts(t *testing.T) {
//	storage, ids := getPopulatedWalletStorage(t)
//	stores.TestingListingAccounts(storage, ids, t)
//}