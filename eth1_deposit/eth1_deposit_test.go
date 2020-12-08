package eth1deposit

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
)

type dummyAccount struct {
	priv *bls.SecretKey
}

func (a *dummyAccount) ID() uuid.UUID                       { return uuid.New() }
func (a *dummyAccount) WalletID() uuid.UUID                 { return uuid.New() }
func (a *dummyAccount) Name() string                        { return "" }
func (a *dummyAccount) PublicKey() *bls.PublicKey           { return a.priv.GetPublicKey() }
func (a *dummyAccount) Path() string                        { return "" }
func (a *dummyAccount) Sign(data []byte) (*bls.Sign, error) { return a.priv.SignByte(data), nil }
func (a *dummyAccount) SetContext(ctx *core.WalletContext)  {}

func _ignoreErr(a []byte, err error) []byte {
	return a
}

// tested against eth2.0-deposit-cli V1.1.0
// Mnemonic: sphere attract wide clown fire balcony dance maple sphere seat design dentist eye orbit diet apart noise cinnamon wealth magic inject witness dress divorce
func TestMainetDepositData(t *testing.T) {
	tests := []struct {
		testname                      string
		validatorPrivKey              []byte
		withdrawalPubKey              []byte
		expectedWithdrawalCredentials []byte
		expectedSig                   []byte
		expectedRoot                  []byte
	}{
		{
			validatorPrivKey:              _ignoreErr(hex.DecodeString("175db1c5411459893301c3f2ebe740e5da07db8f17c2df4fa0be6d31a48a4f79")),
			withdrawalPubKey:              _ignoreErr(hex.DecodeString("8d176708b908f288cc0e9d43f75674e73c0db94026822c5ce2c3e0f9e773c9ee95fdba824302f1208c225b0ed2d54154")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16")),
			expectedSig:                   _ignoreErr(hex.DecodeString("8ab63bb2ef45d5fe4b5ba3b6aa2db122db350c05846b6ffc1415c603ba998226599a21aa65a8cb55c1b888767bdac2b51901d34cde41003c689b8c125fc67d3abd2527ccaf1390c13c3fc65a7422de8a7e29ae8e9736321606172c7b3bf6de36")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("76139d2c8d8e87a4737ce7acbf97ce8980732921550c5443a8754635c11296d3")),
		},
	}

	require.NoError(t, core.InitBLS())

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			val, err := core.NewHDKeyFromPrivateKey(test.validatorPrivKey, "")
			require.NoError(t, err)

			// create data
			depositData, root, err := DepositData(
				val,
				test.withdrawalPubKey,
				core.MainNetwork,
				MaxEffectiveBalanceInGwei,
			)
			require.NoError(t, err)
			require.Equal(t, val.PublicKey().Serialize(), depositData.PublicKey)
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

// tested against eth2.0-deposit-cli V1.1.0
// Mnemonic: sphere attract wide clown fire balcony dance maple sphere seat design dentist eye orbit diet apart noise cinnamon wealth magic inject witness dress divorce
func TestPyrmontDepositData(t *testing.T) {
	tests := []struct {
		testname                      string
		validatorPrivKey              []byte
		withdrawalPubKey              []byte
		expectedWithdrawalCredentials []byte
		expectedSig                   []byte
		expectedRoot                  []byte
	}{
		{
			validatorPrivKey:              _ignoreErr(hex.DecodeString("175db1c5411459893301c3f2ebe740e5da07db8f17c2df4fa0be6d31a48a4f79")),
			withdrawalPubKey:              _ignoreErr(hex.DecodeString("8d176708b908f288cc0e9d43f75674e73c0db94026822c5ce2c3e0f9e773c9ee95fdba824302f1208c225b0ed2d54154")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16")),
			expectedSig:                   _ignoreErr(hex.DecodeString("ab4c9da6a20f385da7a8beb8ac8f58d691d83cd31ba807dbb6de631f5b4c5b1e82e811e41422ccbcd16ef5cb370e50af093dd58ebbc1575b5ed0395ab94538bf0a938f75ec683d4e2e1090c6f1e79a85d771781c3d72a3718451684360e43241")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("06175367bbd24e1966fb7b4299d1a6b0fc107c4385872fc9e4f956e5ffcb61dc")),
		},
	}

	require.NoError(t, core.InitBLS())

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			val, err := core.NewHDKeyFromPrivateKey(test.validatorPrivKey, "")
			require.NoError(t, err)

			// create data
			depositData, root, err := DepositData(
				val,
				test.withdrawalPubKey,
				core.PyrmontNetwork,
				MaxEffectiveBalanceInGwei,
			)
			require.NoError(t, err)
			require.Equal(t, val.PublicKey().Serialize(), depositData.PublicKey)
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

func TestUnsupportedNetwork(t *testing.T) {
	depositData, root, err := DepositData(
		nil,
		make([]byte, 48),
		core.Network("not_supported"),
		MaxEffectiveBalanceInGwei,
	)
	require.EqualError(t, err, "Network not_supported is not supported")
	require.Nil(t, depositData)
	require.EqualValues(t, root, [32]byte{})
}
