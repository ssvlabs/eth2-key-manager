package wallet_hd

import (
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshaling(t *testing.T) {
	tests := []struct{
		id uuid.UUID
		testName string
		walletType core.WalletType
		name string
		indexMapper map[string]uuid.UUID
	}{
		{
			testName:"simple wallet, no account",
			id:uuid.New(),
			walletType:core.HDWallet,
			name: "wallet1",
			indexMapper:map[string]uuid.UUID{},
		},
		{
			testName:"simple wallet with accounts",
			id:uuid.New(),
			walletType:core.HDWallet,
			name: "wallet1",
			indexMapper:map[string]uuid.UUID{
				"account1" : uuid.New(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			w := &HDWallet{
				walletType:test.walletType,
				name:test.name,
				id:test.id,
				indexMapper:test.indexMapper,
			}

			// marshal
			byts,err := json.Marshal(w)
			if err != nil {
				t.Error(err)
				return
			}
			//unmarshal
			w1 := &HDWallet{}
			err = json.Unmarshal(byts,w1)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t,w.id,w1.id)
			require.Equal(t,w.name,w1.name)
			require.Equal(t,w.walletType,w1.walletType)
			for k := range w.indexMapper {
				v := w.indexMapper[k]
				require.Equal(t,v,w1.indexMapper[k])
			}
		})
	}

}