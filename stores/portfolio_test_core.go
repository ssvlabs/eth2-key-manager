package stores

import (
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"os"
	"testing"
)

func portfolio(storage core.Storage) (core.Portfolio,error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	options := &KeyVault.PortfolioOptions{}
	options.SetStorage(storage)
	options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"))
	return KeyVault.NewKeyVault(options)
}

func TestingNonExistingPortfolio(storage core.Storage, t *testing.T) {
	w, err := storage.OpenPortfolio()
	if err != nil {
		t.Error("returned an error for a non existing wallet, should not return an error but rather a nil wallet")
		return
	}

	if w != nil {
		t.Error("returned a wallet for a non existing uuid")
	}
}


func TestingPortfolioStorage(storage core.Storage, t *testing.T) {
	tests := []struct{
		name string
		encryptor types.Encryptor
		password []byte
		error
	}{
		{
			name:"serialization and fetching",
		},
		{
			name:"serialization and fetching with encryptor",
			encryptor: keystorev4.New(),
			password: []byte("password"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p,err := portfolio(storage)
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

			err = storage.SavePortfolio(p)
			if err != nil {
				if test.error != nil {
					require.Equal(t,test.error.Error(),err.Error())
				} else {
					t.Error(err)
				}
				return
			}

			// fetch wallet by id
			fetched, err := storage.OpenPortfolio()
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
			require.Equal(t,p.ID(),fetched.ID())
		})
	}

	// reset
	storage.SetEncryptor(nil,nil)
}