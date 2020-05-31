package stores

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
	"strings"
	"testing"
)

func TestingAccountAtNonExistingWallet(storage wtypes.Store, _ []uuid.UUID, t *testing.T) {
	walletId := uuid.New()
	err := storage.StoreAccount(walletId,uuid.New(),[]byte("account 1"))
	if err == nil {
		t.Error(fmt.Errorf("no error was thrown"))
	}

	expectedErr := fmt.Sprintf("wallet not found")
	if strings.Compare(err.Error(),expectedErr) != 0 {
		t.Error(fmt.Errorf("expeced error: %s but received: %s",expectedErr,err.Error()))
	}
}

func TestingAddingAccountsToWallet(storage wtypes.Store, ids []uuid.UUID, t *testing.T) {
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

func TestingFetchingNonExistingAccount(storage wtypes.Store, ids []uuid.UUID, t *testing.T) {
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

func TestingListingAccounts(storage wtypes.Store, ids []uuid.UUID, t *testing.T) {
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