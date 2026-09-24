package inmemory

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	eth2keymanager "github.com/ssvlabs/eth2-key-manager"
	"github.com/ssvlabs/eth2-key-manager/core"
	encryptor2 "github.com/ssvlabs/eth2-key-manager/encryptor"
	"github.com/ssvlabs/eth2-key-manager/encryptor/keystorev4"
	"github.com/ssvlabs/eth2-key-manager/wallets/hd"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func encryptor() encryptor2.Encryptor {
	return keystorev4.New()
}

func TestStoringWithEncryption(t *testing.T) {
	storage := getStorage()
	tests := []struct {
		testName string
		password []byte
		secret   []byte
		err      error
	}{
		{
			testName: "secret smaller than 32 bytes, should error",
			password: []byte("12345"),
			secret:   []byte("some seed"),
			err:      errors.New("secret can be only 32 bytes (not 9 bytes)"),
		},
		{
			testName: "secret longer than 32 bytes, should error",
			password: []byte("12345"),
			secret:   []byte("i am much longer than 32 bytes of data believe me people!"),
			err:      errors.New("secret can be only 32 bytes (not 57 bytes)"),
		},
		{
			testName: "secret exactly 32 bytes",
			password: []byte("12345"),
			secret:   []byte("i am exactly 32 bytes, pass me!!"),
		},
		{
			testName: "password empty string",
			password: []byte(""),
			secret:   []byte("i am exactly 32 bytes, pass me!!"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// set encryptor
			storage.SetEncryptor(encryptor(), test.password)

			w := hd.NewWallet(&core.WalletContext{Storage: storage})

			err := storage.SaveWallet(w)
			require.NoError(t, err)

			w1, err := storage.OpenWallet()
			require.NoError(t, err)
			require.NotNil(t, w1)
			require.Equal(t, w.ID(), w1.ID())
		})
	}
}

func getPopulatedWalletStorage(t *testing.T) (core.Storage, []core.ValidatorAccount, error) {
	require.NoError(t, core.InitBLS())
	store := getStorage()

	// seed
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")

	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	vault, err := eth2keymanager.NewKeyVault(options)
	if err != nil {
		return nil, nil, err
	}

	wallet, err := vault.Wallet()
	if err != nil {
		return nil, nil, err
	}

	a1, err := wallet.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, nil, err
	}
	a2, err := wallet.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, nil, err
	}
	a3, err := wallet.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, nil, err
	}
	a4, err := wallet.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, nil, err
	}

	return store, []core.ValidatorAccount{a1, a2, a3, a4}, nil
}

func TestOpeningAccount(t *testing.T) {
	storage, accounts, err := getPopulatedWalletStorage(t)
	require.NoError(t, err)
	a1, err := storage.OpenAccount(accounts[0].ID())
	require.NoError(t, err)
	require.Equal(t, accounts[0].ID().String(), a1.ID().String())
	require.Equal(t, accounts[0].ValidatorPublicKey(), a1.ValidatorPublicKey())
	require.Equal(t, accounts[0].WithdrawalPublicKey(), a1.WithdrawalPublicKey())
	require.Equal(t, accounts[0].Name(), a1.Name())
}

func TestAddingAccountsToWallet(t *testing.T) {
	storage, accounts, err := getPopulatedWalletStorage(t)
	require.NoError(t, err)
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

func TestFetchingNonExistingAccount(t *testing.T) {
	storage, _, err := getPopulatedWalletStorage(t)
	require.NoError(t, err)
	t.Run("testing", func(t *testing.T) {
		// fetch non existing account
		_, err := storage.OpenAccount(uuid.New())
		require.NoError(t, err)
	})
}

func TestListingAccounts(t *testing.T) {
	storage, _, err := getPopulatedWalletStorage(t)
	require.NoError(t, err)
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

func getStorage() core.Storage {
	return NewInMemStore(core.MainNetwork)
}

func keyVault(storage core.Storage) (*eth2keymanager.KeyVault, error) {
	if err := core.InitBLS(); err != nil {
		os.Exit(1)
	}

	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(storage)
	return eth2keymanager.NewKeyVault(options)
}

func TestOpeningAccounts(t *testing.T) {
	storage := getStorage()
	kv, err := keyVault(storage)
	require.NoError(t, err)

	wallet, err := kv.Wallet()
	require.NoError(t, err)

	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")

	for i := 0; i < 10; i++ {
		testName := fmt.Sprintf("adding and fetching account: %d", i)
		t.Run(testName, func(t *testing.T) {
			// create
			a, err := wallet.CreateValidatorAccount(seed, nil)
			require.NoError(t, err)

			// open
			a1, err := wallet.AccountByPublicKey(hex.EncodeToString(a.ValidatorPublicKey()))
			require.NoError(t, err)

			a2, err := wallet.AccountByID(a.ID())
			require.NoError(t, err)

			// verify
			for _, fetchedAccount := range []core.ValidatorAccount{a1, a2} {
				require.Equal(t, a.ID().String(), fetchedAccount.ID().String())
				require.Equal(t, a.Name(), fetchedAccount.Name())
				require.Equal(t, a.ValidatorPublicKey(), fetchedAccount.ValidatorPublicKey())
				require.Equal(t, a.WithdrawalPublicKey(), fetchedAccount.WithdrawalPublicKey())
			}
		})
	}
}

func TestNonExistingWallet(t *testing.T) {
	storage := getStorage()
	w, err := storage.OpenWallet()
	require.NotNil(t, err)
	require.EqualError(t, err, "wallet not found")
	require.Nil(t, w)
}

func TestWalletStorage(t *testing.T) {
	storage := getStorage()
	tests := []struct {
		name       string
		walletName string
		encryptor  encryptor2.Encryptor
		password   []byte
		error
	}{
		{
			name:       "serialization and fetching",
			walletName: "test1",
		},
		{
			name:       "serialization and fetching with encryptor",
			walletName: "test2",
			encryptor:  keystorev4.New(),
			password:   []byte("password"),
		},
	}

	kv, err := keyVault(storage)
	require.NoError(t, err)

	wallet, err := kv.Wallet()
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// set encryptor
			if test.encryptor != nil {
				storage.SetEncryptor(test.encryptor, test.password)
			} else {
				storage.SetEncryptor(nil, nil)
			}

			err = storage.SaveWallet(wallet)
			if err != nil {
				if test.error != nil {
					require.Equal(t, test.Error(), err.Error())
				} else {
					t.Error(err)
				}
				return
			}

			// fetch wallet by id
			fetched, err := storage.OpenWallet()
			if err != nil {
				if test.error != nil {
					require.Equal(t, test.Error(), err.Error())
				} else {
					t.Error(err)
				}
				return
			}

			require.NotNil(t, fetched)
			require.NoError(t, test.error)

			// assert
			require.Equal(t, wallet.ID(), fetched.ID())
			require.Equal(t, wallet.Type(), fetched.Type())
		})
	}

	// reset
	storage.SetEncryptor(nil, nil)
}
