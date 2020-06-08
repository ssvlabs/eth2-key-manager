package stores

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"testing"
)

type dummyEncryptor struct {
}
func (enc *dummyEncryptor)Name() string { return "" }
func (enc *dummyEncryptor)Version() uint {return 0}
func (enc *dummyEncryptor) Encrypt(data []byte, key []byte) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data":string(data),
	},nil
}
func (enc *dummyEncryptor) Decrypt(data map[string]interface{}, key []byte) ([]byte, error) {
	return []byte(data["data"].(string)),nil
}


func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func TestingWithdrawalAccount(storage core.PortfolioStorage, t *testing.T) {
	portfolio,err := portfolio(storage)
	if err != nil {
		t.Error(err)
		return
	}
	w,err := portfolio.CreateWallet("test")
	if err != nil {
		t.Error(err)
		return
	}

	a,err := w.GetWithdrawalAccount()
	if err != nil {
		t.Error(err)
		return
	}

	require.Equal(t,"wallet_withdrawal_key_unique",a.Name())
	require.Equal(t,"m/12381/3600/0/0",a.Path())
	require.Equal(t,"ad05971d7f95df8e9181d3d4fb298643493e32df84735ba98d63ddeb4c18492c191973c6dfb01d08ed4f93e7d8247933",hex.EncodeToString(a.PublicKey().Marshal()))
}

func TestingOpenAccounts(storage core.PortfolioStorage, t *testing.T) {
	portfolio,err := portfolio(storage)
	if err != nil {
		t.Error(err)
		return
	}
	w,err := portfolio.CreateWallet("test")
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0 ; i < 10 ; i ++ {
		testName := fmt.Sprintf("adding and fetching account: %d", i)
		t.Run(testName, func(t *testing.T) {
			// create
			accountName := fmt.Sprintf("%d",i)
			a,err := w.CreateValidatorAccount(accountName)
			if err != nil {
				t.Error(err)
				return
			}

			// open
			a1,err := w.AccountByName(accountName)
			if err != nil {
				t.Error(err)
				return
			}
			a2,err := w.AccountByID(a.ID())
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

func TestingNonExistingWallet(storage core.PortfolioStorage, t *testing.T) {
	uid := uuid.New()
	w, err := storage.OpenWallet(uid)
	if err != nil {
		t.Error("returned an error for a non existing wallet, should not return an error but rather a nil wallet")
		return
	}

	if w != nil {
		t.Error("returned a wallet for a non existing uuid")
	}
}

func TestingWalletListing(storage core.PortfolioStorage, t *testing.T) {
	wallets := []struct{
		name string
		walletName string
	}{
		{
			name:"add wallet 1",
			walletName:"1",
		},
		{
			name:"add wallet 2",
			walletName:"2",
		},
		{
			name:"add wallet 3",
			walletName:"3",
		},
		{
			name:"add wallet 4",
			walletName:"4",
		},
	}

	ids := make([]uuid.UUID,4)
	portfolio,err := portfolio(storage)
	if err != nil {
		t.Error(err)
		return
	}

	// store
	t.Run("storing", func(t *testing.T) {
		for i, test := range wallets {
			w,err := portfolio.CreateWallet(test.walletName)
			if err != nil {
				t.Error(err)
				return
			}

			err = storage.SaveWallet(w)
			if err != nil {
				t.Error(err)
			}

			ids[i] = w.ID()
		}

	})

	// util
	findFunc := func(slice []uuid.UUID, val uuid.UUID) (int, bool) {
		for i, item := range slice {
			if item == val {
				return i, true
			}
		}
		return -1, false
	}

	matched := 0
	// fetch all
	t.Run("fetching", func(t *testing.T) {
		wallets, err := storage.ListWallets()
		if err != nil {
			t.Error(err)
			return
		}

		for _,w := range wallets {
			if _,res := findFunc(ids,w.ID()); res {
				matched ++
			}
		}

		require.Equal(t,len(ids),matched, "not all saved wallets found")
	})
}

func TestingWalletStorage(storage core.PortfolioStorage, t *testing.T) {
	tests := []struct{
		name string
		walletName string
		encryptor types.Encryptor
		password []byte
		error
	}{
		{
			name:"serialization and fetching",
			walletName:"test",
		},
		{
			name:"serialization and fetching with encryptor",
			walletName:"test",
			encryptor: &dummyEncryptor{},
			password: []byte("password"),
		},
		{
			name:"serialization and fetching with encryptor (no password)",
			walletName:"test",
			encryptor: &dummyEncryptor{},
			error: fmt.Errorf("can't encrypt, missing password"),
		},
	}

	portfolio,err := portfolio(storage)
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w,err := portfolio.CreateWallet(test.walletName)
			if err != nil {
				t.Error(err)
				return
			}

			// set encryptor
			if test.encryptor != nil {
				storage.SetEncryptor(test.encryptor,test.password)
			} else {
				storage.SetEncryptor(nil,nil)
			}

			err = storage.SaveWallet(w)
			if err != nil {
				if test.error != nil {
					require.Equal(t,test.error.Error(),err.Error())
				} else {
					t.Error(err)
				}
				return
			}

			// fetch wallet by id
			fetched, err := storage.OpenWallet(w.ID())
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
			require.Equal(t,w.ID(),fetched.ID())
			require.Equal(t,w.Name(),fetched.Name())
			require.Equal(t,w.Type(),fetched.Type())
		})
	}

	// reset
	storage.SetEncryptor(nil,nil)
}

//func TestingUpdatingWallet(storage core.PortfolioStorage, t *testing.T) {
//	id := uuid.New()
//
//	// new wallet
//	err := storage.StoreWallet(id,"wallet",[]byte("wallet"))
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	// add accounts
//	err = storage.StoreAccount(id,uuid.New(),[]byte("account 1"))
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	err = storage.StoreAccount(id,uuid.New(),[]byte("account 2"))
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	// update just wallet data
//	err = storage.StoreWallet(id,"wallet",[]byte("wallet updated"))
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	// verify wallet data
//	value, err := storage.RetrieveWalletByID(id)
//	if err != nil {
//		t.Error(err)
//	}
//	if bytes.Compare([]byte("wallet updated"),value) != 0 {
//		t.Error(fmt.Errorf("did not retrieve the same data, required: %s, received: %s","wallet updated", string(value)))
//		return
//	}
//	// verify accounts
//	expectedAccountCnt := 0
//	for _ = range storage.RetrieveAccounts(id) {
//		expectedAccountCnt ++
//	}
//	if expectedAccountCnt != 2 {
//		t.Error(fmt.Errorf("expected %d accounts, recieved: %d",2,expectedAccountCnt))
//	}
//}