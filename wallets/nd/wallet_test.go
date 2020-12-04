package nd

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/dummy"
)

func storage() core.Storage {
	return &dummy.DummyStorage{}
}

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func TestWalletMarshaling(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		testName    string
		walletType  core.WalletType
		indexMapper map[string]uuid.UUID
		seed        []byte
		path        string
	}{
		{
			testName:    "simple wallet, no account",
			id:          uuid.New(),
			walletType:  core.HDWallet,
			indexMapper: map[string]uuid.UUID{},
			seed:        _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path:        "/0/0",
		},
		{
			testName:   "simple wallet with accounts",
			id:         uuid.New(),
			walletType: core.HDWallet,
			indexMapper: map[string]uuid.UUID{
				"ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279": uuid.New(),
			},
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// setup storage
			storage := storage()

			w := &Wallet{
				walletType:  test.walletType,
				id:          test.id,
				indexMapper: test.indexMapper,
				//key:key,
			}

			// marshal
			byts, err := json.Marshal(w)
			require.NoError(t, err)

			//unmarshal
			w1 := &Wallet{context: &core.WalletContext{Storage: storage}}
			err = json.Unmarshal(byts, w1)
			require.NoError(t, err)

			require.Equal(t, w.id, w1.id)
			require.Equal(t, w.walletType, w1.walletType)
			for k := range w.indexMapper {
				v := w.indexMapper[k]
				require.Equal(t, v, w1.indexMapper[k])
			}
		})
	}

}
