package hashicorp

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
	"strings"
	"testing"
)

func getPopulatedWalletStorage(t *testing.T) (wtypes.Store,[]uuid.UUID) {
	store := getWalletStorage()

	ids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	if err := store.StoreWallet(ids[0],"1", []byte("wallet 1")); err != nil { t.Error(err) }
	if err := store.StoreWallet(ids[1],"2", []byte("wallet 2")); err != nil { t.Error(err) }
	if err := store.StoreWallet(ids[2],"3", []byte("wallet 3")); err != nil { t.Error(err) }
	if err := store.StoreWallet(ids[3],"4", []byte("wallet 4")); err != nil { t.Error(err) }

	return store,ids
}

func TestAccoountAtNonExistingWallet(t *testing.T) {
	storage, _ := getPopulatedWalletStorage(t)

	walletId := uuid.New()
	err := storage.StoreAccount(walletId,uuid.New(),[]byte("account 1"))
	if err == nil {
		t.Error(fmt.Errorf("no error was thrown"))
	}

	expectedErr := fmt.Sprintf("could not retrieve wallet id: %s",walletId.String())
	if strings.Compare(err.Error(),expectedErr) != 0 {
		t.Error(fmt.Errorf("expeced error: %s but received: %s",expectedErr,err.Error()))
	}
}

func TestAddingAccountsToWallet(t *testing.T) {
	storage, ids := getPopulatedWalletStorage(t)

	for i := 0 ; i < 10 ; i++ {
		accountId := uuid.New()
		testname := fmt.Sprintf("adding account %s",accountId.String())
		t.Run(testname, func(t *testing.T) {
			err := storage.StoreAccount(ids[0],accountId,[]byte(accountId.String()))
			if err != nil {
				t.Error(err)
				return
			}

			// verify account was added
			val,err := storage.RetrieveAccount(ids[0],accountId)
			if err != nil {
				t.Error(err)
			}
			if res := bytes.Compare(val,[]byte(accountId.String())); res != 0 {
				t.Error(fmt.Errorf("could not fetch stored account %s",accountId.String()))
			}
		})
	}
}

func TestFetchingNonExistingAccount(t *testing.T) {
	storage, ids := getPopulatedWalletStorage(t)
	t.Run("testing", func(t *testing.T) {
		err := storage.StoreAccount(ids[0],uuid.New(),[]byte("test"))
		if err != nil {
			t.Error(err)
			return
		}

		// fetch non existing account
		_,err = storage.RetrieveAccount(ids[0],uuid.New())
		if err == nil {
			t.Error(fmt.Errorf("did not return error when fetching non existing account"))
		}
	})
}

func TestListingAccounts(t *testing.T) {
	storage, ids := getPopulatedWalletStorage(t)

	accounts := map[string]bool{}
	// add accounts
	for i := 0 ; i < 10 ; i++ {
		accountId := uuid.New()
		err := storage.StoreAccount(ids[0],accountId,[]byte(accountId.String()))
		if err != nil {
			t.Error(err)
			return
		}
		accounts[accountId.String()] = false
	}



	// verify listing
	for a := range storage.RetrieveAccounts(ids[0]) {
		accountid := string(a)
		accounts[accountid] = true
	}
	for k,v := range accounts {
		t.Run(k, func(t *testing.T) {
			if v != true {
				t.Error(fmt.Errorf("account %s not fetched",k))
				return
			}
		})
	}

}