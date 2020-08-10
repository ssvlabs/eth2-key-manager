package eth1_deposit

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"os"
	"testing"
)

type dummyAccount struct {
	priv *e2types.BLSPrivateKey
}

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func newDummyAccount(privKey []byte) *dummyAccount {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	k, err := e2types.BLSPrivateKeyFromBytes(privKey)
	if err != nil {
		return nil
	}
	return &dummyAccount{priv: k}
}
func (a *dummyAccount) ID() uuid.UUID                               { return uuid.New() }
func (a *dummyAccount) WalletID() uuid.UUID                         { return uuid.New() }
func (a *dummyAccount) Name() string                                { return "" }
func (a *dummyAccount) PublicKey() e2types.PublicKey                { return a.priv.PublicKey() }
func (a *dummyAccount) Path() string                                { return "" }
func (a *dummyAccount) Sign(data []byte) (e2types.Signature, error) { return a.priv.Sign(data), nil }
func (a *dummyAccount) SetContext(ctx *core.WalletContext)          {}

func _ignoreErr(a []byte, err error) []byte {
	return a
}

func TestDepositData(t *testing.T) {
	tests := []struct {
		testname                      string
		seed                          []byte
		validatorPrivKey              []byte
		withdrawalPrivKey             []byte
		expectedWithdrawalCredentials []byte
		expectedSig                   []byte
		expectedRoot                  []byte
	}{
		{
			seed:                          _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			validatorPrivKey:              _ignoreErr(hex.DecodeString("2ca74ab9c41ec9cacdd5f1543633104d272d24a64aeaf9bf200a79653b446333")),
			withdrawalPrivKey:             _ignoreErr(hex.DecodeString("03beaaad608afa088f459be520985b4f67796c0ffda51b2c6301a50bb74d131b")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("0080e08282d1bf1182645774610d621f77f03ae4f9e692202cf896b1914eed28")),
			expectedSig:                   _ignoreErr(hex.DecodeString("acfdf9926446f691f88926a6ce2f4534d00b10dce3f44cfb7ec9138424b56f40db884962484ea20b704591fad57a5b8513ebf13ed11235d4fa23dcccd33ab7312c864a750e1bda4d91cd1cd4ad20b5ac6e96aacdce34b775af6a7ddeaf44341e")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("0a65fff5cac97ec7a2a170b81aa07a0cc04c7a3e198d55c3e5f254115465cb08")),
		},
	}

	e2types.InitBLS()

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			seed := test.seed
			options := &KeyVault.KeyVaultOptions{}
			options.SetStorage(in_memory.NewInMemStore())
			options.SetSeed(seed)
			options.SetEncryptor(keystorev4.New())
			options.SetPassword("password")
			kv, err := KeyVault.NewKeyVault(options)
			require.NoError(t, err)
			wallet, err := kv.Wallet()
			require.NoError(t, err)

			// create and fetch validator account
			val, err := wallet.CreateValidatorAccount(seed, "val1")
			require.NoError(t, err)
			withd, err := core.NewHDKeyFromPrivateKey(test.withdrawalPrivKey, "")
			require.NoError(t, err)

			// create data
			depositData, root, err := DepositData(val, withd, MaxEffectiveBalanceInGwei)
			require.NoError(t, err)

			require.Equal(t, val.ValidatorPublicKey().Marshal(), depositData.PublicKey)
			require.True(t, bytes.Equal(test.expectedWithdrawalCredentials, depositData.WithdrawalCredentials))
			require.Equal(t, MaxEffectiveBalanceInGwei, depositData.Amount)
			require.True(t, bytes.Equal(test.expectedRoot, root[:]))
			require.True(t, bytes.Equal(test.expectedSig, depositData.Signature))

			fmt.Printf("pubkey: %s\n", hex.EncodeToString(depositData.PublicKey))
			fmt.Printf("WithdrawalCredentials: %s\n", hex.EncodeToString(depositData.WithdrawalCredentials))
			fmt.Printf("Amount: %d\n", depositData.Amount)
			fmt.Printf("root: %s\n", hex.EncodeToString(root[:]))
			fmt.Printf("sig: %s\n", hex.EncodeToString(depositData.Signature))
		})
	}
}
