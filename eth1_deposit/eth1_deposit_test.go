package eth1_deposit

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"os"
	"testing"
)

type dummyAccount struct {
	priv *e2types.BLSPrivateKey
}
func newDummyAccount(privKey []byte) *dummyAccount {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	k,err := e2types.BLSPrivateKeyFromBytes(privKey)
	if err != nil {
		return nil
	}
	return &dummyAccount{priv:k}
}
func (a *dummyAccount) ID() uuid.UUID {return uuid.New()}
func (a *dummyAccount) WalletID() uuid.UUID                        {return uuid.New()}
func (a *dummyAccount) Type() core.AccountType                     {return core.ValidatorAccount }
func (a *dummyAccount) Name() string                               {return ""}
func (a *dummyAccount) PublicKey() e2types.PublicKey {return a.priv.PublicKey()}
func (a *dummyAccount) Path() string {return ""}
func (a *dummyAccount) Sign(data []byte) (e2types.Signature,error) {return a.priv.Sign(data),nil}
func (a *dummyAccount) SetContext(ctx *core.WalletContext)         {}

func _ignoreErr(a []byte, err error) []byte {
	return a
}

func TestDepositData(t *testing.T) {
	tests := []struct{
		testname string
		validatorPrivKey []byte
		withdrawalPrivKey []byte
		expectedWithdrawalCredentials []byte
		expectedSig []byte
		expectedRoot []byte
	}{
		{
			validatorPrivKey: _ignoreErr(hex.DecodeString("3811c38debe9275248b9d873fe886c7e428ded40a95eda3fc00e9b147000d8ec")),
			withdrawalPrivKey: _ignoreErr(hex.DecodeString("64d5b3638d3bc555073d6ab3a3652b117bb1c3715ae4f17c0b2fa26e88a9f660")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("00c76a029adcac82fe161b34f44de3c8c94182ffe75bf29a938691ebfd66bf6b")),
			expectedSig: _ignoreErr(hex.DecodeString("88ff6c5a44b85db96b684cee772506489ae388838fe4d13435bf415de23ce14a9b4f254dd1f456cffbb581d87f4a6ce806f559e8d1afa28cdbde84a5fba6526e9f948ddde7166d8ba8218478e5e681833492d61a7b49d11ced0718ac317218df")),
			expectedRoot: _ignoreErr(hex.DecodeString("6255505dc4c2ba5828cc6ad8f47bd122f02d8c840fc1aa81abd817f3971c2d79")),
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			val := newDummyAccount(test.validatorPrivKey)
			withd := newDummyAccount(test.withdrawalPrivKey)

			// create data
			depositData,root,err := DepositData(val,withd, MaxEffectiveBalanceInGwei)
			require.NoError(t,err)

			require.Equal(t,val.PublicKey().Marshal(),depositData.PublicKey)
			require.True(t,bytes.Equal(test.expectedWithdrawalCredentials,depositData.WithdrawalCredentials))
			require.Equal(t, MaxEffectiveBalanceInGwei,depositData.Amount)
			require.True(t,bytes.Equal(test.expectedRoot,root[:]))
			require.True(t,bytes.Equal(test.expectedSig,depositData.Signature))


			fmt.Printf("pubkey: %s\n",hex.EncodeToString(depositData.PublicKey))
			fmt.Printf("WithdrawalCredentials: %s\n",hex.EncodeToString(depositData.WithdrawalCredentials))
			fmt.Printf("Amount: %d\n",depositData.Amount)
			fmt.Printf("root: %s\n",hex.EncodeToString(root[:]))
			fmt.Printf("sig: %s\n",hex.EncodeToString(depositData.Signature))
		})
	}
}

