package stores

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func keyVault(storage core.Storage) (*eth2keymanager.KeyVault, error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(storage)
	options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"))
	return eth2keymanager.NewKeyVault(options)
}

func TestingOpenAccounts(storage core.Storage, t *testing.T) {
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
			if err != nil {
				t.Error(err)
				return
			}

			// open
			a1, err := wallet.AccountByPublicKey(hex.EncodeToString(a.ValidatorPublicKey().Marshal()))
			if err != nil {
				t.Error(err)
				return
			}
			a2, err := wallet.AccountByID(a.ID())
			if err != nil {
				t.Error(err)
				return
			}

			// verify
			for _, fetchedAccount := range []core.ValidatorAccount{a1, a2} {
				require.Equal(t, a.ID().String(), fetchedAccount.ID().String())
				require.Equal(t, a.Name(), fetchedAccount.Name())
				require.Equal(t, a.ValidatorPublicKey().Marshal(), fetchedAccount.ValidatorPublicKey().Marshal())
				require.Equal(t, a.WithdrawalPublicKey().Marshal(), fetchedAccount.WithdrawalPublicKey().Marshal())
			}
		})
	}

}

func TestingNonExistingWallet(storage core.Storage, t *testing.T) {
	w, err := storage.OpenWallet()
	require.NotNil(t, err)
	require.EqualError(t, err, "wallet not found")
	require.Nil(t, w)
}

func TestingWalletStorage(storage core.Storage, t *testing.T) {
	tests := []struct {
		name       string
		walletName string
		encryptor  types.Encryptor
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
					require.Equal(t, test.error.Error(), err.Error())
				} else {
					t.Error(err)
				}
				return
			}

			// fetch wallet by id
			fetched, err := storage.OpenWallet()
			if err != nil {
				if test.error != nil {
					require.Equal(t, test.error.Error(), err.Error())
				} else {
					t.Error(err)
				}
				return
			}
			if fetched == nil {
				t.Errorf("wallet could not be fetched by id")
				return
			}

			if test.error != nil {
				t.Errorf("expected error: %s", test.error.Error())
				return
			}

			// assert
			require.Equal(t, wallet.ID(), fetched.ID())
			require.Equal(t, wallet.Type(), fetched.Type())
		})
	}

	// reset
	storage.SetEncryptor(nil, nil)
}
