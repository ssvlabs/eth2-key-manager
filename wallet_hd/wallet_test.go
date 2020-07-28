package wallet_hd

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"math/big"
	"os"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func inmemStorage() *in_memory.InMemStore {
	return in_memory.NewInMemStore()
}

func key(seed []byte) (*core.MasterDerivableKey, error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	return core.MasterKeyFromSeed(seed)
}

func TestAccountDerivation(t *testing.T) {
	// create wallet
	storage := inmemStorage()
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	w := &HDWallet{
		id:          uuid.New(),
		indexMapper: make(map[string]uuid.UUID),
		//key:key,
		context: &core.WalletContext{
			Storage: storage,
		},
	}

	tests := []struct {
		testName              string
		accountName           string
		expectedValidationKey *big.Int
		expectedWithdrawalKey *big.Int
	}{
		{
			testName:              "account 0",
			accountName:           "account 0",
			expectedValidationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			expectedWithdrawalKey: _bigInt("51023953445614749789943419502694339066585011438324100967164633618358653841358"),
		},
		{
			testName:              "account 1",
			accountName:           "account 1",
			expectedValidationKey: _bigInt("22295543756806915021696580341385697374834805500065673451566881420621123341007"),
			expectedWithdrawalKey: _bigInt("19211358943475501217006127435996279333633291783393046900803879394346849035913"),
		},
		{
			testName:              "account 2",
			accountName:           "account 2",
			expectedValidationKey: _bigInt("43442610958028244518598118443083802862055489983359071059993155323547905350874"),
			expectedWithdrawalKey: _bigInt("23909010000215292098635609623453075881965979294359727509549907878193079139650"),
		},
		{
			testName:              "account 3",
			accountName:           "account 3",
			expectedValidationKey: _bigInt("4448413729621370906608934836012354998323947125552823486758689486871003717293"),
			expectedWithdrawalKey: _bigInt("37328169013635701905066231905928437636499300152882617419715404470232404314068"),
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			account, err := w.CreateValidatorAccount(seed, test.accountName)
			require.NoError(t, err)

			val, err := e2types.BLSPrivateKeyFromBytes(test.expectedValidationKey.Bytes())
			require.NoError(t, err)
			with, err := e2types.BLSPrivateKeyFromBytes(test.expectedWithdrawalKey.Bytes())
			require.NoError(t, err)

			require.Equal(t, val.PublicKey().Marshal(), account.ValidatorPublicKey().Marshal())
			require.Equal(t, with.PublicKey().Marshal(), account.WithdrawalPublicKey().Marshal())
		})
	}
}

func TestCreateAccounts(t *testing.T) {
	tests := []struct {
		testName        string
		newAccounttName string
		expectedErr     string
	}{
		{
			testName:        "Add new account",
			newAccounttName: "account1",
			expectedErr:     "",
		},
		{
			testName:        "Add duplicate account, should error",
			newAccounttName: "account1",
			expectedErr:     "account account1 already exists",
		},
		{
			testName:        "Add account with no name, should error",
			newAccounttName: "account1",
			expectedErr:     "account name is empty",
		},
	}

	// create key and wallet
	storage := inmemStorage()
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")

	w := &HDWallet{
		id:          uuid.New(),
		indexMapper: make(map[string]uuid.UUID),
		//key:key,
		context: &core.WalletContext{
			Storage: storage,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_, err := w.CreateValidatorAccount(seed, test.newAccounttName)
			if test.expectedErr != "" {
				require.Errorf(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
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
				"account1": uuid.New(),
			},
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// setup storage
			storage := inmemStorage()

			w := &HDWallet{
				walletType:  test.walletType,
				id:          test.id,
				indexMapper: test.indexMapper,
				//key:key,
			}

			// marshal
			byts, err := json.Marshal(w)
			if err != nil {
				t.Error(err)
				return
			}
			//unmarshal
			w1 := &HDWallet{context: &core.WalletContext{Storage: storage}}
			err = json.Unmarshal(byts, w1)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t, w.id, w1.id)
			require.Equal(t, w.walletType, w1.walletType)
			for k := range w.indexMapper {
				v := w.indexMapper[k]
				require.Equal(t, v, w1.indexMapper[k])
			}
		})
	}

}
