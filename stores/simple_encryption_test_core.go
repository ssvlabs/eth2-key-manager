package stores

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"testing"
)

func encryptor() types.Encryptor {
	return keystorev4.New()
}

func TestingPortfolioStorageWithEncryption(storage core.Storage, t *testing.T) {
	tests := []struct{
		testName string
		password []byte
		secret []byte
		err error
	}{
		{
			testName:"secret smaller than 32 bytes, should error",
			password:[]byte("12345"),
			secret: []byte("some seed"),
			err:fmt.Errorf("secret can be only 32 bytes (not 9 bytes)"),
		},
		{
			testName:"secret longer than 32 bytes, should error",
			password:[]byte("12345"),
			secret: []byte("i am much longer than 32 bytes of data beleive me people!"),
			err: fmt.Errorf("secret can be only 32 bytes (not 57 bytes)"),
		},
		{
			testName:"secret exactly 32 bytes",
			password:[]byte("12345"),
			secret: []byte("i am exactly 32 bytes, pass me!!"),
		},
		{
			testName:"password empty string",
			password:[]byte(""),
			secret: []byte("i am exactly 32 bytes, pass me!!"),
		},
	}

	for _,test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// set encryptor
			storage.SetEncryptor(encryptor(),test.password)

			// store a secret seed
			err := storage.SecurelySavePortfolioSeed(test.secret)
			if test.err != nil {
				require.EqualError(t,err,test.err.Error())
				return
			} else if err != nil {
				t.Error(err)
				return
			}

			// fetch and compare
			ret,err := storage.SecurelyFetchPortfolioSeed()
			if err != nil {
				t.Error(err)
				return
			}
			require.NotNil(t,ret)
			require.Equal(t,test.secret,ret[:len(test.secret)])
		})
	}
}
