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
			expectedSig:                   _ignoreErr(hex.DecodeString("ad3a400c3aadeaf5f734ba88511bcd27872a4561080126b967e46d3742c9b62c62ff93503a166c37382868c46816fd58083db53731d0c5413dc2801a2308ffb35d18997779bf1af01cd76489ad42d91bb67211dd02723b728f8a8a08c3307a77")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("fb8defd3efaa1c73967bc4624e5f6ad548ffef348223a713f1118dc585a77fca")),
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
			withdrawalPubKey:              _ignoreErr(hex.DecodeString("8d176708b908f288cc0e9d43f75674e73c0db94026822c5ce2c3e0f9e773c9ee95fdba824302f1208c225b0ed2d54154")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16")),
			expectedSig:                   _ignoreErr(hex.DecodeString("aacdb59866b8092f004233799a41ca488b606c937689f1e09905476577269b6819ebf25ad0fc44e54799cc57852a5e19126b667fd9fd0df73d7e8d9f24203eb26ad920a0fdcc9e60f2e300fd0a3caf64b1fa19c59bf5dfb84fd14176948a92d2")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("b29c77d193afa3b6caadd22a845d39d047aaef991927e031e6fbbb4b6995b5f4")),
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
