package eth1deposit

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
)

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
			require.Equal(t, val.PublicKey().SerializeToHexStr(), strings.TrimPrefix(depositData.PublicKey.String(), "0x"))
			require.Equal(t, test.expectedWithdrawalCredentials, depositData.WithdrawalCredentials)
			require.Equal(t, MaxEffectiveBalanceInGwei, depositData.Amount)
			require.Equal(t, test.expectedRoot, root[:], hex.EncodeToString(root[:]))
			require.Equal(t, hex.EncodeToString(test.expectedSig), strings.TrimPrefix(depositData.Signature.String(), "0x"))

			fmt.Printf("pubkey: %s\n", hex.EncodeToString(depositData.PublicKey[:]))
			fmt.Printf("WithdrawalCredentials: %s\n", hex.EncodeToString(depositData.WithdrawalCredentials))
			fmt.Printf("Amount: %d\n", depositData.Amount)
			fmt.Printf("root: %s\n", hex.EncodeToString(root[:]))
			fmt.Printf("sig: %s\n", hex.EncodeToString(depositData.Signature[:]))
		})
	}
}

// tested against eth2.0-deposit-cli V1.1.0
// Mnemonic: sphere attract wide clown fire balcony dance maple sphere seat design dentist eye orbit diet apart noise cinnamon wealth magic inject witness dress divorce
func TestPraterDepositData(t *testing.T) {
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
			withdrawalPubKey:              _ignoreErr(hex.DecodeString("b3d50de8d77299da8d830de1edfb34d3ce03c1941846e73870bb33f6de7b8a01383f6b32f55a1d038a4ddcb21a765194")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("006029659d86cf9f19d53030273372c84b1912d0633cb15381a75cb92850f03a")),
			expectedSig:                   _ignoreErr(hex.DecodeString("a2bcc9d2ac82062cb9806b761e8e8d405963620b8f5356fa70fe543812bf07c3031546482c737401ba1dec01d5690d0600c900ebe7dca5699e804ff4441ed4e25789b389bcdc69c6f4dc25ef40e5694f6de7723bda359c5c2a54e05ae90290ca")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("d243130779e16b4352bb8d2c80765334b4a7bdd4bc42356b37e42380dc47dac5")),
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
				core.PraterNetwork,
				MaxEffectiveBalanceInGwei,
			)

			require.NoError(t, err)
			require.Equal(t, val.PublicKey().SerializeToHexStr(), strings.TrimPrefix(depositData.PublicKey.String(), "0x"))
			require.Equal(t, test.expectedWithdrawalCredentials, depositData.WithdrawalCredentials)
			require.Equal(t, MaxEffectiveBalanceInGwei, depositData.Amount)
			require.Equal(t, test.expectedRoot, root[:], hex.EncodeToString(root[:]))
			require.Equal(t, hex.EncodeToString(test.expectedSig), strings.TrimPrefix(depositData.Signature.String(), "0x"))

			fmt.Printf("pubkey: %s\n", hex.EncodeToString(depositData.PublicKey[:]))
			fmt.Printf("WithdrawalCredentials: %s\n", hex.EncodeToString(depositData.WithdrawalCredentials))
			fmt.Printf("Amount: %d\n", depositData.Amount)
			fmt.Printf("root: %s\n", hex.EncodeToString(root[:]))
			fmt.Printf("sig: %s\n", hex.EncodeToString(depositData.Signature[:]))
		})
	}
}

func TestUnsupportedNetwork(t *testing.T) {
	depositData, root, err := DepositData(
		nil,
		make([]byte, 48),
		"not_supported",
		MaxEffectiveBalanceInGwei,
	)
	require.EqualError(t, err, "Network not_supported is not supported")
	require.Nil(t, depositData)
	require.EqualValues(t, root, [32]byte{})
}
