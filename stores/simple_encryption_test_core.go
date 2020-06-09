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

func encryptor() types.Encryptor {
	return keystorev4.New()
}

func TestingPortfolioStorageWithEncryption(storage core.Storage, t *testing.T) {
	tests := []struct{
		testName string
		password []byte
	}{
		{
			testName:"password empty string",
			password:[]byte(""),
		},
	}

	for _,test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// set encryptor
			storage.SetEncryptor(encryptor(),test.password)

			// create portfolio
			if err := e2types.InitBLS(); err != nil {
				os.Exit(1)
			}

			options := &KeyVault.PortfolioOptions{}
			options.SetStorage(storage)
			v,err := KeyVault.NewKeyVault(options)
			if err != nil {
				t.Error(err)
				return
			}

			// save portfolio
			storage.SavePortfolio(v)

			// fetch portfolio
			v1,err := storage.OpenPortfolio()
			if err != nil {
				t.Error(err)
				return
			}
			if v1 == nil {
				t.Errorf("could not open portfolio")
				return
			}

			require.Equal(t,v.ID().String(),v1.ID().String())
		})
	}
}
