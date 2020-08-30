package eth1_deposit

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth-key-manager/core"
)

type dummyAccount struct {
	priv *e2types.BLSPrivateKey
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
			expectedSig:                   _ignoreErr(hex.DecodeString("91c8dddcbd882d409e5b0eca4bf096decf349edf4bd118cf32bcee73acd698a436f8a4ff1fe16f041df963d712bed4d916dd00a8bd13f3087369fe199086d46a1c788633d8fda2ba4ace3e2ece56dc66b947fb6d83941bd2568cd2fcfd03c298")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("7155ceb606f46c5102003121bd19935590a58fb28ec2efe345c66db622008874")),
		},
	}

	e2types.InitBLS()

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			val, err := core.NewHDKeyFromPrivateKey(test.validatorPrivKey, "")
			require.NoError(t, err)

			// create data
			depositData, root, err := DepositData(val, test.withdrawalPubKey, MaxEffectiveBalanceInGwei)
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
