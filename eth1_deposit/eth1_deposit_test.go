package eth1deposit

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
	types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

func _ignoreErr(a []byte, err error) []byte {
	return a
}

// tested against eth2.0-deposit-cli V1.1.0
// Mnemonic: sphere attract wide clown fire balcony dance maple sphere seat design dentist eye orbit diet apart noise cinnamon wealth magic inject witness dress divorce
func TestMainetDepositData(t *testing.T) {
	tests := []struct {
		network                       core.Network
		testname                      string
		validatorPrivKey              []byte
		withdrawalPubKey              []byte
		expectedWithdrawalCredentials []byte
		expectedSig                   []byte
		expectedRoot                  []byte
	}{
		{
			network:                       core.MainNetwork,
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
				test.network,
				MaxEffectiveBalanceInGwei,
			)
			VerifyOperation(t, depositData, test.network)

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
		network                       core.Network
		testname                      string
		validatorPrivKey              []byte
		withdrawalPubKey              []byte
		expectedWithdrawalCredentials []byte
		expectedSig                   []byte
		expectedRoot                  []byte
	}{
		{
			network:                       core.PraterNetwork,
			validatorPrivKey:              _ignoreErr(hex.DecodeString("175db1c5411459893301c3f2ebe740e5da07db8f17c2df4fa0be6d31a48a4f79")),
			withdrawalPubKey:              _ignoreErr(hex.DecodeString("8d176708b908f288cc0e9d43f75674e73c0db94026822c5ce2c3e0f9e773c9ee95fdba824302f1208c225b0ed2d54154")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16")),
			expectedSig:                   _ignoreErr(hex.DecodeString("a88d0fd588836c5756ec7f2fe2bc8b6fc5723d018c8d31c8f42b239ac6cf7c2f9ae129caafaebb5f2f25e7821678b41819bc24f6eeebe0d8196cea13581f72ac501f3e7e9e4bc596e6a545ac109fb2ff1d7eb03923454dc5258718b43427a757")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("aa940a26af67a676bcd807b0fd3f39aadbfc6862e380e115051683e1fccc0171")),
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
				test.network,
				MaxEffectiveBalanceInGwei,
			)
			VerifyOperation(t, depositData, test.network)

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
func TestHoleskyDepositData(t *testing.T) {
	tests := []struct {
		network                       core.Network
		testname                      string
		validatorPrivKey              []byte
		withdrawalPubKey              []byte
		expectedWithdrawalCredentials []byte
		expectedSig                   []byte
		expectedRoot                  []byte
	}{
		{
			network:                       core.HoleskyNetwork,
			validatorPrivKey:              _ignoreErr(hex.DecodeString("175db1c5411459893301c3f2ebe740e5da07db8f17c2df4fa0be6d31a48a4f79")),
			withdrawalPubKey:              _ignoreErr(hex.DecodeString("8d176708b908f288cc0e9d43f75674e73c0db94026822c5ce2c3e0f9e773c9ee95fdba824302f1208c225b0ed2d54154")),
			expectedWithdrawalCredentials: _ignoreErr(hex.DecodeString("005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16")),
			expectedSig:                   _ignoreErr(hex.DecodeString("836bccc57ceb05353119814a025d8a83a271d6724d1eb760d1c806e9de15a919f389cd6235e6a6b1bda4cfd3c236882c1858bcf4b3141d3a3fba73c158ce59d28adcf2e67dbf05dc00d944a47cfd8db08a8de7a145f2f4c6888714be77b410e2")),
			expectedRoot:                  _ignoreErr(hex.DecodeString("75e81e6fde731d5a2f5360af3baca7d1cb599ed10288df3bd7988e9f7ad8c929")),
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
				test.network,
				MaxEffectiveBalanceInGwei,
			)
			VerifyOperation(t, depositData, test.network)

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

func VerifyOperation(t *testing.T, depositData *phase0.DepositData, network core.Network) {
	depositMessage := &phase0.DepositMessage{
		WithdrawalCredentials: depositData.WithdrawalCredentials,
		Amount:                depositData.Amount,
	}
	copy(depositMessage.PublicKey[:], depositData.PublicKey[:])

	depositMsgRoot, err := depositMessage.HashTreeRoot()
	require.NoError(t, err)
	require.NotNil(t, depositMsgRoot)

	sigBytes := make([]byte, len(depositData.Signature))
	copy(sigBytes, depositData.Signature[:])
	sig, err := types.BLSSignatureFromBytes(sigBytes)
	require.NoError(t, err)
	require.NotNil(t, sig)

	container := &phase0.SigningData{
		ObjectRoot: depositMsgRoot,
	}

	genesisForkVersion := network.GenesisForkVersion()
	domain, err := types.ComputeDomain(types.DomainDeposit, genesisForkVersion[:], types.ZeroGenesisValidatorsRoot)
	require.NoError(t, err)
	copy(container.Domain[:], domain[:])
	signingRoot, err := container.HashTreeRoot()
	require.NoError(t, err)
	require.NotNil(t, signingRoot)

	var pubkeyBytes [48]byte
	copy(pubkeyBytes[:], depositData.PublicKey[:])

	pubkey, err := types.BLSPublicKeyFromBytes(pubkeyBytes[:])
	require.NoError(t, err)
	require.NotNil(t, pubkey)
	require.True(t, sig.Verify(signingRoot[:], pubkey))
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
