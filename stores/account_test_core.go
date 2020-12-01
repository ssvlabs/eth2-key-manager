package stores

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
)

func TestingOpeningAccount(storage core.Storage, account core.ValidatorAccount, t *testing.T) {
	a1, err := storage.OpenAccount(account.ID())
	require.NoError(t, err)
	require.Equal(t, account.ID().String(), a1.ID().String())
	require.Equal(t, account.ValidatorPublicKey(), a1.ValidatorPublicKey())
	require.Equal(t, account.WithdrawalPublicKey(), a1.WithdrawalPublicKey())
	require.Equal(t, account.Name(), a1.Name())
}

func TestingSavingAccounts(storage core.Storage, accounts []core.ValidatorAccount, t *testing.T) {
	for _, account := range accounts {
		t.Run(fmt.Sprintf("adding account %s", account.Name()), func(t *testing.T) {
			err := storage.SaveAccount(account)
			require.NoError(t, err)

			// verify account was added
			val, err := storage.OpenAccount(account.ID())
			require.NoError(t, err)
			require.Equal(t, account.ID(), val.ID())
			require.Equal(t, account.Name(), val.Name())
			require.Equal(t, account.ValidatorPublicKey(), val.ValidatorPublicKey())
			require.Equal(t, account.WithdrawalPublicKey(), val.WithdrawalPublicKey())
		})
	}
}

func TestingFetchingNonExistingAccount(storage core.Storage, t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		// fetch non existing account
		_, err := storage.OpenAccount(uuid.New())
		require.NoError(t, err)
	})
}

func TestingListingAccounts(storage core.Storage, t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	// create keyvault and wallet
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(storage)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)

	wallet, err := vault.Wallet()
	require.NoError(t, err)

	// create accounts
	accounts := map[string]bool{}
	for i := 0; i < 10; i++ {
		account, err := wallet.CreateValidatorAccount(seed, nil)
		require.NoError(t, err)

		accounts[account.ID().String()] = false
	}

	// verify listing
	fetched, err := storage.ListAccounts()
	require.NoError(t, err)

	for _, a := range fetched {
		accounts[a.ID().String()] = true
	}
	for k, v := range accounts {
		t.Run(k, func(t *testing.T) {
			require.True(t, v)
		})
	}
}
