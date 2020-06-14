package KeyVault

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
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

func inmemStorage() *in_memory.InMemStore {
	return in_memory.NewInMemStore()
}

func key(seed []byte, relativePath string, storage core.Storage) (*core.DerivableKey,error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	key,err := core.BaseKeyFromSeed(seed, storage)
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
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0",
		},
		{
			testName:"simple portfolio with wallets",
			id:uuid.New(),
			simpleSigner: false,
			indexMapper:map[string]uuid.UUID{
				"wallet1" : uuid.New(),
			},
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0",
		},
		{
			testName:"simple portfolio without key path",
			id:uuid.New(),
			simpleSigner: false,
			indexMapper:map[string]uuid.UUID{
				"wallet1" : uuid.New(),
			},
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// setup storage
			storage := inmemStorage()
			err := storage.SecurelySavePortfolioSeed(test.seed)
			if err != nil {
				t.Error(err)
				return
			}

			// create key and vault
			key,err := key(test.seed,test.path, storage)
			if err != nil {
				t.Error(err)
				return
			}
			w := &KeyVault{
				id:test.id,
				indexMapper:test.indexMapper,
				key:key,
				Context:&core.PortfolioContext{Storage:storage},
			}

			// marshal
			byts,err := json.Marshal(w)
			if err != nil {
				t.Error(err)
				return
			}
			//unmarshal
			w1 := &KeyVault{Context:&core.PortfolioContext{Storage:storage},}
			err = json.Unmarshal(byts,w1)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t,w.id,w1.id) // id
			for k := range w.indexMapper { // index mapper
				v := w.indexMapper[k]
				require.Equal(t,v,w1.indexMapper[k])
			}
			require.Equal(t,w.key.PublicKey().Marshal(),w1.key.PublicKey().Marshal()) // key
		})
	}

}
