package eth1_deposit

import (
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
			expectedSig:                   _ignoreErr(hex.DecodeString("b922154b5e1ab3302e0bba98e3eee2f94e8ee246622264e9fd6364530be1e9c94ce76648780b118cfac5741a62abf05b061009c2a41d8a459f2accab91564b5b67b38c53e62037a85cbf366ba63e0b88073e93821e8c9de1c87749f2db925aef")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("bc610cc4fe56d60c64e6665dae24dee70968a6e23cf52ee1da90d99adedcf250")),
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
			require.Equal(t, test.expectedWithdrawalCredentials, depositData.WithdrawalCredentials)
			require.Equal(t, MaxEffectiveBalanceInGwei, depositData.Amount)
			require.Equal(t, test.expectedRoot, root[:], hex.EncodeToString(root[:]))
			require.Equal(t, test.expectedSig, depositData.Signature, hex.EncodeToString(depositData.Signature))

			fmt.Printf("pubkey: %s\n", hex.EncodeToString(depositData.PublicKey))
			fmt.Printf("WithdrawalCredentials: %s\n", hex.EncodeToString(depositData.WithdrawalCredentials))
			fmt.Printf("Amount: %d\n", depositData.Amount)
			fmt.Printf("root: %s\n", hex.EncodeToString(root[:]))
			fmt.Printf("sig: %s\n", hex.EncodeToString(depositData.Signature))
		})
	}
}
