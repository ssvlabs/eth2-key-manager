package wallet_hd

import (
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountMarshaling(t *testing.T) {
	tests := []struct{
		id uuid.UUID
		testName string
		accounttType core.AccountType
		name string
		seed []byte
		path string
	}{
		{
			testName:"simple account",
			id:uuid.New(),
			accounttType:core.ValidatorAccount,
			name: "account1",
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
			a := &HDAccount{
				accountType:test.accounttType,
				name:test.name,
				id:test.id,
				key:key,
			}

			// marshal
			byts,err := json.Marshal(a)
			if err != nil {
				t.Error(err)
				return
			}
			//unmarshal
			a1 := &HDAccount{}
			err = json.Unmarshal(byts,a1)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t,a.id,a1.id)
			require.Equal(t,a.name,a1.name)
			require.Equal(t,a.accountType,a1.accountType)
			require.Equal(t,a.key.Key.Marshal(),a1.key.Key.Marshal())
		})
	}
}