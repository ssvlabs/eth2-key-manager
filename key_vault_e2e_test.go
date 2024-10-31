package eth2keymanager

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/eth2-key-manager/core"
	"github.com/ssvlabs/eth2-key-manager/encryptor/keystorev4"
	"github.com/ssvlabs/eth2-key-manager/stores/inmemory"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func inmemStorage() *inmemory.InMemStore {
	return inmemory.NewInMemStore(core.MainNetwork)
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
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")
			v, err := NewKeyVault(options)
			require.NoError(t, err)

			// test common tests
			testVault(t, v, seed)

			wallet, err := v.Wallet()
			require.NoError(t, err)

			// test specific derivation
			account, err := wallet.AccountByPublicKey("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")
			require.NoError(t, err)
			require.NotNil(t, account)

			require.Equal(t, _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"), account.ValidatorPublicKey())
			require.Equal(t, _byteArray("a0b9324da8a8a696c53950e984de25b299c123d17bab972eca1ac2c674964c9f817047bc6048ef0705d7ec6fae6d5da6"), account.WithdrawalPublicKey())
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
			importedVault, err := NewKeyVault(options)
			require.NoError(t, err)
			// test common tests
			testVault(t, importedVault, test.seed) // this will create some wallets and accounts

			// open vault
			v, err := OpenKeyVault(options)
			require.NoError(t, err)

			wallet, err := v.Wallet()
			require.NoError(t, err)

			// test specific derivation
			account, err := wallet.AccountByPublicKey("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")
			require.NoError(t, err)
			require.NotNil(t, account)

			require.Equal(t, _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"), account.ValidatorPublicKey())
			require.Equal(t, _byteArray("a0b9324da8a8a696c53950e984de25b299c123d17bab972eca1ac2c674964c9f817047bc6048ef0705d7ec6fae6d5da6"), account.WithdrawalPublicKey())
			require.Equal(t, importedVault.walletID, v.walletID)
		})
	}
}

func testVault(t *testing.T, v *KeyVault, seed []byte) {
	wallet, err := v.Wallet()
	require.NoError(t, err)

	// create and fetch validator account
	val, err := wallet.CreateValidatorAccount(seed, nil)
	require.NoError(t, err)
	val1, err := wallet.AccountByPublicKey(hex.EncodeToString(val.ValidatorPublicKey()))
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
