package stores

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
	wtypes "github.com/wealdtech/go-eth2-wallet-types"
	"strings"
	"testing"
)

func getStorage() logical.Storage {
	return &logical.InmemStorage{}
}

func getWalletStorage() wtypes.Store {
	return NewHashicorpVaultStore(getStorage(),context.Background())
}

func TestNonExistingWallet(t *testing.T) {
	storage := getWalletStorage()
	uid := uuid.New()
	_, err := storage.RetrieveWalletByID(uid)
	if err == nil {
		fmt.Errorf("returned a non error for a non existing wallet")
	}

	expectedErr := fmt.Sprintf("could not retrieve wallet id: %s",uid.String())
	if err.Error() != expectedErr {
		t.Error(fmt.Errorf("errors not match, required: %s, received: %s",expectedErr,err.Error()))
	}
}

func TestMultiWalletStorage(t *testing.T) {
	storage := getWalletStorage()

	wallets := []struct{
		name string
		walletId uuid.UUID
		walletName string
		data []byte
		error string
	}{
		{
			name:"1",
			walletId:uuid.New(),
			walletName:"1",
			data: []byte("1"),
		},
		{
			name:"2",
			walletId:uuid.New(),
			walletName:"2",
			data: []byte("2"),
		},
		{
			name:"3",
			walletId:uuid.New(),
			walletName:"3",
			data: []byte("3"),
		},
		{
			name:"4",
			walletId:uuid.New(),
			walletName:"4",
			data: []byte("4"),
		},
	}

	// store
	t.Run("storing", func(t *testing.T) {
		for _, test := range wallets {
			err := storage.StoreWallet(test.walletId,test.walletName, test.data)
			if err != nil {
				t.Error(err)
			}
		}

	})

	// fetch all
	t.Run("fetching", func(t *testing.T) {
		walletnames := map[string]bool{"1":false,"2":false,"3":false,"4":false}
		for w := range storage.RetrieveWallets() {
			walletnames[string(w)] = true
		}

		for k,v := range walletnames {
			if v != true {
				t.Error(fmt.Errorf("Wallet %s not fetched",k))
				return
			}
		}
	})
}

func TestWalletStorage(t *testing.T) {
	storage := getWalletStorage()

	tests := []struct{
		name string
		walletId uuid.UUID
		walletName string
		data []byte
		error string
	}{
		{
			name:"simple data",
			walletId:uuid.New(),
			walletName:"test",
			data: []byte("test data"),
		},
		{
			name:"no wallet name",
			walletId:uuid.New(),
			walletName:"",
			data: []byte("test data"),
			error: "wallet name must be provided",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := storage.StoreWallet(test.walletId,test.walletName, test.data)
			if err != nil {
				if len(test.error) != 0 {
					if res := strings.Compare(test.error,err.Error()); res != 0 {
						t.Error(fmt.Errorf("received wrong error, required: %s, received: %s",test.error, err.Error()))
					}
				} else {
					t.Error(err)
				}
				return
			}

			// by value
			value, err := storage.RetrieveWalletByID(test.walletId)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(test.data,value) != 0 {
				t.Error(fmt.Errorf("did not retrieve the same data, required: %s, received: %s",string(test.data), string(value)))
			}

			// by wallet name
			value, err = storage.RetrieveWallet(test.walletName)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(test.data,value) != 0 {
				t.Error(fmt.Errorf("did not retrieve the same data, required: %s, received: %s",string(test.data), string(value)))
			}
		})
	}
}
