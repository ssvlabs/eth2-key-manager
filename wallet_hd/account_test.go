package wallet_hd

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	types "github.com/wealdtech/go-eth2-types/v2"
	"testing"
)

func TestAccountMarshaling(t *testing.T) {
	tests := []struct{
		id uuid.UUID
		testName string
		//accountType core.AccountType
		parentWalletId uuid.UUID
		name string
		seed []byte
		accountIndex string
	}{
		{
			testName:"simple account",
			id:uuid.New(),
			//accountType:core.ValidatorAccount,
			parentWalletId:uuid.New(),
			name: "account1",
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			accountIndex: "0",
		},
	}

	types.InitBLS()

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// setup storage
			storage := inmemStorage()
			err := storage.SecurelySavePortfolioSeed(test.seed)
			require.NoError(t, err)

			// create key and account
			masterKey,err := core.MasterKeyFromSeed(storage)
			require.NoError(t, err)
			validationKey,err := masterKey.Derive(fmt.Sprintf("/%s/0/0", test.accountIndex))
			require.NoError(t, err)
			withdrawalKey,err := masterKey.Derive(fmt.Sprintf("/%s/0", test.accountIndex))
			require.NoError(t, err)
			a := &HDAccount{
				//accountType:test.accountType,
				name:test.name,
				id:test.id,
				validationKey: validationKey,
				withdrawalPubKey:withdrawalKey.PublicKey(),
			}

			// marshal
			byts,err := json.Marshal(a)
			require.NoError(t, err)
			//unmarshal
			a1 := &HDAccount{context:&core.WalletContext{Storage: storage}}
			err = json.Unmarshal(byts,a1)
			require.NoError(t, err)

			require.Equal(t,a.id,a1.id)
			require.Equal(t,a.name,a1.name)
			//require.Equal(t,a.accountType,a1.accountType)
			//require.Equal(t,a.WalletID().String(),a1.WalletID().String())
			require.Equal(t,a.validationKey.PublicKey().Marshal(),a1.validationKey.PublicKey().Marshal())
			require.Equal(t,a.withdrawalPubKey.Marshal(),a1.withdrawalPubKey.Marshal())
		})
	}
}