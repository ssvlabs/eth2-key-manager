package KeyVault

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"os"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func key(seed []byte, relativePath string) (*core.DerivableKey,error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	key,err := core.BaseKeyFromSeed(seed)
	if err != nil {
		return nil,err
	}

	if len(relativePath) >0 {
		return key.Derive(relativePath)
	} else {
		return key,nil
	}
}

func TestMarshaling(t *testing.T) {
	tests := []struct{
		id uuid.UUID
		simpleSigner bool
		testName string
		indexMapper map[string]uuid.UUID
		seed []byte
		path string
	}{
		{
			testName:"simple portfolio, no wallets",
			id:uuid.New(),
			simpleSigner: true,
			indexMapper:map[string]uuid.UUID{},
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path: "/0/0",
		},
		{
			testName:"simple portfolio with wallets",
			id:uuid.New(),
			simpleSigner: false,
			indexMapper:map[string]uuid.UUID{
				"wallet1" : uuid.New(),
			},
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
			path: "/0/0",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			key,err := key(test.seed,test.path)
			if err != nil {
				t.Error(err)
				return
			}
			w := &KeyVault{
				id:test.id,
				enableSimpleSigner: test.simpleSigner,
				indexMapper:test.indexMapper,
				key:key,
			}

			// marshal
			byts,err := json.Marshal(w)
			if err != nil {
				t.Error(err)
				return
			}
			//unmarshal
			w1 := &KeyVault{}
			err = json.Unmarshal(byts,w1)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t,w.id,w1.id)
			require.Equal(t,w.enableSimpleSigner,w1.enableSimpleSigner)
			for k := range w.indexMapper {
				v := w.indexMapper[k]
				require.Equal(t,v,w1.indexMapper[k])
			}
			require.Equal(t,w.key.Key.Marshal(),w1.key.Key.Marshal())
		})
	}

}
