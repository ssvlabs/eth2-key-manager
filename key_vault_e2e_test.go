package eth2keymanager

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
)

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func inmemStorage() *in_memory.InMemStore {
	return in_memory.NewInMemStore(core.MainNetwork)
}

func TestNoKeyVault(t *testing.T) {
	tests := []struct {
		testname string
		storage  core.Storage
	}{
		{
			testname: "In-memory storage",
			storage:  inmemStorage(),
		},
		{
			testname: "Hashicorp Vault storage",
			storage:  inmemStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			options := &KeyVaultOptions{}
			options.SetStorage(test.storage)

			kv, err := OpenKeyVault(options)
			require.NotNil(t, err)
			require.EqualError(t, err, "wallet not found")
			require.Nil(t, kv)
		})
	}
}

func TestNewKeyVault(t *testing.T) {
	tests := []struct {
		testname string
		storage  core.Storage
	}{
		{
			testname: "In-memory storage",
			storage:  inmemStorage(),
		},
		{
			testname: "Hashicorp Vault storage",
			storage:  inmemStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			// setup vault
			options := &KeyVaultOptions{}
			options.SetStorage(test.storage)
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")
			v, err := NewKeyVault(options)
			require.NoError(t, err)

			// generate new seed
			seed, err := core.GenerateNewEntropy()
			require.NoError(t, err)

			testVault(t, v, seed)
		})
	}
}

func TestImportKeyVault(t *testing.T) {
	tests := []struct {
		testname string
		seed     []byte
		storage  core.Storage
	}{
		{
			testname: "In-memory storage",
			seed:     _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:  inmemStorage(),
		},
		{
			testname: "Hashicorp Vault storage",
			seed:     _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:  inmemStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			seed := test.seed
			options := &KeyVaultOptions{}
			options.SetStorage(test.storage)
			options.SetSeed(seed)
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")
			v, err := NewKeyVault(options)
			require.NoError(t, err)

			// test common tests
			testVault(t, v, seed)

			wallet, err := v.Wallet()
			require.NoError(t, err)

			// test specific derivation
			account, err := wallet.AccountByPublicKey("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279")
			require.NoError(t, err)
			require.NotNil(t, account)

			expectedValKey, err := e2types.BLSPrivateKeyFromBytes(_bigInt("16278447180917815188301017385774271592438483452880235255024605821259671216398").Bytes())
			require.NoError(t, err)
			expectedWithdrawalKey, err := e2types.BLSPrivateKeyFromBytes(_bigInt("26551663876804375121305275007227133452639447817512639855729535822239507627836").Bytes())
			require.NoError(t, err)

			require.Equal(t, expectedValKey.PublicKey().Marshal(), account.ValidatorPublicKey().Marshal())
			require.Equal(t, expectedWithdrawalKey.PublicKey().Marshal(), account.WithdrawalPublicKey().Marshal())
		})
	}
}

func TestOpenKeyVault(t *testing.T) {
	tests := []struct {
		testname string
		seed     []byte
		storage  core.Storage
	}{
		{
			testname: "In-memory storage",
			seed:     _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:  inmemStorage(),
		},
		{
			testname: "Hashicorp Vault storage",
			seed:     _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:  inmemStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			// options
			storage := test.storage
			storage.SetEncryptor(keystorev4.New(), []byte("password"))
			options := &KeyVaultOptions{}
			options.SetStorage(storage)
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")

			// import keyvault
			options.SetSeed(test.seed)
			importedVault, err := NewKeyVault(options)
			// test common tests
			testVault(t, importedVault, test.seed) // this will create some wallets and accounts

			// open vault
			options.SetSeed(nil) // important
			v, err := OpenKeyVault(options)
			require.NoError(t, err)

			wallet, err := v.Wallet()
			require.NoError(t, err)

			// test specific derivation
			account, err := wallet.AccountByPublicKey("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279")
			require.NoError(t, err)
			require.NotNil(t, account)

			expectedValKey, err := e2types.BLSPrivateKeyFromBytes(_bigInt("16278447180917815188301017385774271592438483452880235255024605821259671216398").Bytes())
			require.NoError(t, err)
			expectedWithdrawalKey, err := e2types.BLSPrivateKeyFromBytes(_bigInt("26551663876804375121305275007227133452639447817512639855729535822239507627836").Bytes())
			require.NoError(t, err)

			require.Equal(t, expectedValKey.PublicKey().Marshal(), account.ValidatorPublicKey().Marshal())
			require.Equal(t, expectedWithdrawalKey.PublicKey().Marshal(), account.WithdrawalPublicKey().Marshal())
			require.Equal(t, importedVault.walletId, v.walletId)
		})
	}
}

func testVault(t *testing.T, v *KeyVault, seed []byte) {
	wallet, err := v.Wallet()
	require.NoError(t, err)

	// create and fetch validator account
	val, err := wallet.CreateValidatorAccount(seed, nil)
	require.NoError(t, err)
	val1, err := wallet.AccountByPublicKey(hex.EncodeToString(val.ValidatorPublicKey().Marshal()))
	require.NoError(t, err)
	val2, err := wallet.AccountByID(val.ID())
	require.NoError(t, err)
	require.NotNil(t, val1)
	require.NotNil(t, val2)
	require.Equal(t, val.ID().String(), val1.ID().String())
	require.Equal(t, val.ID().String(), val2.ID().String())
	require.Equal(t, val.Name(), val1.Name())
	require.Equal(t, val.Name(), val2.Name())
}
