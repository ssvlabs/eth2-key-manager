package KeyVault

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshaling(t *testing.T) {
	tests := []struct{
		id uuid.UUID
		simpleSigner bool
		testName string
		indexMapper map[string]uuid.UUID
	}{
		{
			testName:"simple portfolio, no wallets",
			id:uuid.New(),
			simpleSigner: true,
			indexMapper:map[string]uuid.UUID{},
		},
		{
			testName:"simple portfolio with wallets",
			id:uuid.New(),
			simpleSigner: false,
			indexMapper:map[string]uuid.UUID{
				"wallet1" : uuid.New(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			w := &KeyVault{
				id:test.id,
				enableSimpleSigner: test.simpleSigner,
				indexMapper:test.indexMapper,
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
		})
	}

}
