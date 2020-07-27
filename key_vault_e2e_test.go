package KeyVault

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"math/big"
	"testing"
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
	return in_memory.NewInMemStore()
}

func TestNewKeyVault(t *testing.T) {
	tests := []struct{
		testname string
		storage core.Storage
	}{
		{
			testname:"In-memory storage",
			storage:inmemStorage(),
		},
		{
			testname:"Hashicorp Vault storage",
			storage:inmemStorage(),
		},
	}

	for _,test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			// setup vault
			options := &WalletOptions{}
			options.SetStorage(test.storage)
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")
			v,err := NewKeyVault(options)
			require.NoError(t,err)

			testVault(t,v)
		})
	}
}

func TestImportKeyVault(t *testing.T) {
	tests := []struct{
		testname string
		seed []byte
		storage core.Storage
	}{
		{
			testname:"In-memory storage",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:inmemStorage(),
		},
		{
			testname:"Hashicorp Vault storage",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:inmemStorage(),
		},
	}

	for _,test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			seed := test.seed
			options := &WalletOptions{}
			options.SetStorage(test.storage)
			options.SetSeed(seed)
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")
			v,err := ImportKeyVault(options)
			require.NoError(t,err)

			// test common tests
			testVault(t,v)

			wallet,err := v.Wallet()
			require.NoError(t,err)

			// test specific derivation
			account,err := wallet.AccountByName("val1")
			require.NoError(t,err)
			require.NotNil(t,account)

			expectedValKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218").Bytes())
			require.NoError(t,err)
			expectedWithdrawalKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("51023953445614749789943419502694339066585011438324100967164633618358653841358").Bytes())
			require.NoError(t,err)

			assert.Equal(t,expectedValKey.PublicKey().Marshal(),account.ValidatorPublicKey().Marshal())
			assert.Equal(t,expectedWithdrawalKey.PublicKey().Marshal(), account.WithdrawalPublicKey().Marshal())
		})
	}
}

func TestOpenKeyVault(t *testing.T) {
	tests := []struct{
		testname string
		seed []byte
		storage core.Storage
	}{
		{
			testname:"In-memory storage",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:inmemStorage(),
		},
		{
			testname:"Hashicorp Vault storage",
			seed: _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			storage:inmemStorage(),
		},
	}

	for _,test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			// options
			storage := test.storage
			storage.SetEncryptor(keystorev4.New(), []byte("password"))
			options := &WalletOptions{}
			options.SetStorage(storage)
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")

			// import keyvault
			options.SetSeed(test.seed)
			importedVault,err := ImportKeyVault(options)
			// test common tests
			testVault(t,importedVault) // this will create some wallets and accounts

			// open vault
			options.SetSeed(nil) // important
			v,err := OpenKeyVault(options)
			require.NoError(t,err)

			wallet,err := v.Wallet()
			require.NoError(t,err)

			// test specific derivation
			account,err := wallet.AccountByName("val1")
			require.NoError(t,err)
			require.NotNil(t,account)

			expectedValKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218").Bytes())
			require.NoError(t,err)
			expectedWithdrawalKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("51023953445614749789943419502694339066585011438324100967164633618358653841358").Bytes())
			require.NoError(t,err)

			assert.Equal(t,expectedValKey.PublicKey().Marshal(), account.ValidatorPublicKey().Marshal())
			assert.Equal(t,expectedWithdrawalKey.PublicKey().Marshal(), account.WithdrawalPublicKey().Marshal())
		})
	}
}


func testVault(t *testing.T, v *KeyVault) {
	wallet,err := v.Wallet()
	require.NoError(t,err)

	// create and fetch validator account
	val,err := wallet.CreateValidatorAccount("val1")
	require.NoError(t,err)
	val1,err := wallet.AccountByName("val1")
	require.NoError(t,err)
	val2,err := wallet.AccountByID(val.ID())
	require.NoError(t,err)
	require.NotNil(t,val1)
	require.NotNil(t,val2)
	require.Equal(t,val.ID().String(),val1.ID().String())
	require.Equal(t,val.ID().String(),val2.ID().String())
	require.Equal(t,val.Name(),val1.Name())
	require.Equal(t,val.Name(),val2.Name())
}