package stores

import (
	"fmt"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestingOpeningAccount(storage core.Storage, account core.ValidatorAccount, t *testing.T) {
	a1,err := storage.OpenAccount(account.ID())
	if err != nil {
		t.Error(err)
		return
	}
	require.Equal(t,account.ID().String(),a1.ID().String())
	require.Equal(t,account.ValidatorPublicKey().Marshal(),a1.ValidatorPublicKey().Marshal())
	require.Equal(t,account.WithdrawalPublicKey().Marshal(),a1.WithdrawalPublicKey().Marshal())
	require.Equal(t,account.Name(),a1.Name())
}

func TestingSavingAccounts(storage core.Storage, accounts []core.ValidatorAccount, t *testing.T) {
	for _,account := range accounts {
		testname := fmt.Sprintf("adding account %s",account.Name())
		t.Run(testname, func(t *testing.T) {
			err := storage.SaveAccount(account)
			if err != nil {
				t.Error(err)
				return
			}

			// verify account was added
			val,err := storage.OpenAccount(account.ID())
			if err != nil {
				t.Error(err)
			}
			require.Equal(t,account.ID(), val.ID())
			require.Equal(t,account.Name(), val.Name())
			require.Equal(t,account.ValidatorPublicKey().Marshal(), val.ValidatorPublicKey().Marshal())
			require.Equal(t,account.WithdrawalPublicKey().Marshal(), val.WithdrawalPublicKey().Marshal())
		})
	}
}

func TestingFetchingNonExistingAccount(storage core.Storage, t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		// create keyvault and wallet
		options := &KeyVault.WalletOptions{}
		options.SetStorage(storage)
		options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"))

		// fetch non existing account
		_,err := storage.OpenAccount(uuid.New())
		if err != nil {
			t.Error(fmt.Errorf("should not return error for unknwon account, just nil"))
		}
	})
}

func TestingListingAccounts(storage core.Storage, t *testing.T) {
	// create keyvault and wallet
	options := &KeyVault.WalletOptions{}
	options.SetStorage(storage)
	options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"))
	vault,err := KeyVault.NewKeyVault(options)
	require.NoError(t, err)

	wallet, err := vault.Wallet()
	require.NoError(t, err)

	// create accounts
	accounts := map[string]bool{}
	for i := 0 ; i < 10 ; i++ {
		account,err := wallet.CreateValidatorAccount(fmt.Sprintf("%d",i))
		if err != nil {
			t.Error(err)
			return
		}
		accounts[account.ID().String()] = false
	}

	// verify listing
	fetched,err := storage.ListAccounts()
	if err != nil {
		t.Error(err)
		return
	}
	for _,a := range fetched {
		accounts[a.ID().String()] = true
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