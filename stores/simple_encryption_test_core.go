package stores

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/wallet_hd"
)

func encryptor() types.Encryptor {
	return keystorev4.New()
}

func TestingWalletStorageWithEncryption(storage core.Storage, t *testing.T) {
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
			secret:   []byte("i am much longer than 32 bytes of data beleive me people!"),
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

			w := wallet_hd.NewHDWallet(&core.WalletContext{Storage: storage})

			err := storage.SaveWallet(w)
			require.NoError(t, err)

			w1, err := storage.OpenWallet()
			require.NoError(t, err)
			require.NotNil(t, w1)
			require.Equal(t, w.ID(), w1.ID())
		})
	}
}
