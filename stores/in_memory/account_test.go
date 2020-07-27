package in_memory

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func getPopulatedWalletStorage() (core.Storage,[]core.Account,error) {
	store := getStorage()

	options := &KeyVault.PortfolioOptions{}
	options.SetStorage(store)
	options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"))
	vault,err := KeyVault.NewKeyVault(options)
	if err != nil {
		return nil,nil,err
	}

	wallet,err := vault.CreateWallet("test")
	if err != nil {
		return nil,nil,err
	}

	a1,err := wallet.CreateValidatorAccount("1", nil)
	if err != nil {
		return nil,nil,err
	}
	a2,err := wallet.CreateValidatorAccount("2", nil)
	if err != nil {
		return nil,nil,err
	}
	a3,err := wallet.CreateValidatorAccount("3", nil)
	if err != nil {
		return nil,nil,err
	}
	a4,err := wallet.CreateValidatorAccount("4", nil)
	if err != nil {
		return nil,nil,err
	}

	return store,[]core.Account{a1,a2,a3,a4},nil
}

func TestOpeningAccount (t *testing.T) {
	storage, accounts, err := getPopulatedWalletStorage()
	if err != nil {
		t.Error(err)
		return
	}
	stores.TestingOpeningAccount(storage, accounts[0],t)
}

func TestAddingAccountsToWallet(t *testing.T) {
	storage, accounts, err := getPopulatedWalletStorage()
	if err != nil {
		t.Error(err)
		return
	}
	stores.TestingSavingAccounts(storage, accounts, t)
}

func TestFetchingNonExistingAccount(t *testing.T) {
	storage, _, err := getPopulatedWalletStorage()
	if err != nil {
		t.Error(err)
		return
	}
	stores.TestingFetchingNonExistingAccount(storage, t)
}

func TestListingAccounts(t *testing.T) {
	storage, _, err := getPopulatedWalletStorage()
	if err != nil {
		t.Error(err)
		return
	}
	stores.TestingListingAccounts(storage, t)
}