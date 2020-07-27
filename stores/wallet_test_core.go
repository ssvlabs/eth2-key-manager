package stores

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"os"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func keyVault(storage core.Storage) (*KeyVault.KeyVault,error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	options := &KeyVault.WalletOptions{}
	options.SetStorage(storage)
	options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"))
	return KeyVault.ImportKeyVault(options)
}

func TestingWithdrawalAccount(storage core.Storage, t *testing.T) {
	kv,err := keyVault(storage)
	require.NoError(t, err)

	wallet, err := kv.Wallet()
	require.NoError(t, err)

	a,err := wallet.GetWithdrawalAccount()
	if err != nil {
		t.Error(err)
		return
	}

	require.Equal(t,"wallet_withdrawal_key_unique",a.Name())
	require.Equal(t,"m/12381/3600/0/0",a.Path())
	require.Equal(t,"b08033f01d0c71f63a46117915791187e3257a95552b3701fef80124c492eabe1f10795e684055895b887220460e5f24",hex.EncodeToString(a.PublicKey().Marshal()))
}

func TestingOpenAccounts(storage core.Storage, t *testing.T) {
	kv,err := keyVault(storage)
	require.NoError(t, err)

	wallet, err := kv.Wallet()
	require.NoError(t, err)

	for i := 0 ; i < 10 ; i ++ {
		testName := fmt.Sprintf("adding and fetching account: %d", i)
		t.Run(testName, func(t *testing.T) {
			// create
			accountName := fmt.Sprintf("%d",i)
			a,err := wallet.CreateValidatorAccount(accountName)
			if err != nil {
				t.Error(err)
				return
			}

			// open
			a1,err := wallet.AccountByName(accountName)
			if err != nil {
				t.Error(err)
				return
			}
			a2,err := wallet.AccountByID(a.ID())
			if err != nil {
				t.Error(err)
				return
			}

			// verify
			for _,fetchedAccount := range []core.Account{a1,a2} {
				require.Equal(t,a.ID().String(),fetchedAccount.ID().String())
				require.Equal(t,a.Name(),fetchedAccount.Name())
				require.Equal(t,a.PublicKey().Marshal(),fetchedAccount.PublicKey().Marshal())
				require.Equal(t,fmt.Sprintf("m/12381/3600/0/0/%d",i),fetchedAccount.Path())
			}
		})
	}

}

func TestingNonExistingWallet(storage core.Storage, t *testing.T) {
	w, err := storage.OpenWallet()
	if err != nil {
		t.Error("returned an error for a non existing wallet, should not return an error but rather a nil wallet")
		return
	}

	if w != nil {
		t.Error("returned a wallet for a non existing uuid")
	}
}

func TestingWalletStorage(storage core.Storage, t *testing.T) {
	tests := []struct{
		name string
		walletName string
		encryptor types.Encryptor
		password []byte
		error
	}{
		{
			name:"serialization and fetching",
			walletName:"test1",
		},
		{
			name:"serialization and fetching with encryptor",
			walletName:"test2",
			encryptor: keystorev4.New(),
			password: []byte("password"),
		},
	}

	kv,err := keyVault(storage)
	require.NoError(t, err)

	wallet, err := kv.Wallet()
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// set encryptor
			if test.encryptor != nil {
				storage.SetEncryptor(test.encryptor,test.password)
			} else {
				storage.SetEncryptor(nil,nil)
			}

			err = storage.SaveWallet(wallet)
			if err != nil {
				if test.error != nil {
					require.Equal(t,test.error.Error(),err.Error())
				} else {
					t.Error(err)
				}
				return
			}

			// fetch wallet by id
			fetched, err := storage.OpenWallet()
			if err != nil {
				if test.error != nil {
					require.Equal(t,test.error.Error(),err.Error())
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
			require.Equal(t,wallet.ID(),fetched.ID())
			require.Equal(t,wallet.Type(),fetched.Type())
		})
	}

	// reset
	storage.SetEncryptor(nil,nil)
}