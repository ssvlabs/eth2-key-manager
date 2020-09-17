package eth1_deposit

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

type dummyAccount struct {
	priv *e2types.BLSPrivateKey
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
		validatorPrivKey              []byte
		withdrawalPubKey              []byte
		expectedWithdrawalCredentials []byte
		expectedSig                   []byte
		expectedRoot                  []byte
	}{
		{
			validatorPrivKey:              _ignoreErr(hex.DecodeString("23fd464c122d7fa8c9c8e46d710ae478ab920c8c0587e86556aa968191d5210e")),
			withdrawalPubKey:              _ignoreErr(hex.DecodeString("b323537b2867d9f2bae068f93e75a9e2e1c8d594e3696c34dc8010dc403eaeeaf43756a440fc82e1c6f45c6e8348343f")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("00ea056bfaa692b4e12bb1c3f59049dabcfb0b63f427025c718f5e3b81fdb945")),
			expectedSig:                   _ignoreErr(hex.DecodeString("aac3de8d5d1700e2519da9346625273ded81a4250bd1b98c50e6587acac9545a1d1598472823aea29220cbc45ae9062f09791dc22252efdd3a1531964e4d62a59511e0f332fb3cc5ea7fe0831de696f040fe806f9f22bd29db0466047584cb23")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("5b508bbed40a083809e4d0ee74135c7289020e33e2dbad2e69f41772d09f5a63")),
		},
	}

	e2types.InitBLS()

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			val, err := core.NewHDKeyFromPrivateKey(test.validatorPrivKey, "")
			require.NoError(t, err)

			// create data
			depositData, root, err := DepositData(
				val,
				test.withdrawalPubKey,
				core.TestNetwork,
				MaxEffectiveBalanceInGwei,
			)
			require.NoError(t, err)

			require.Equal(t, val.PublicKey().Marshal(), depositData.PublicKey)
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
