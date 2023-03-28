package signer

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	util "github.com/wealdtech/go-eth2-util"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	prot "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/bloxapp/eth2-key-manager/wallets"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func _byteArray32(input string) [32]byte {
	res, _ := hex.DecodeString(input)
	var res32 [32]byte
	copy(res32[:], res)
	return res32
}

func _byteArray96(input string) [96]byte {
	res, _ := hex.DecodeString(input)
	var res96 [96]byte
	copy(res96[:], res)
	return res96
}

// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/master/shared/testutil/block.go#L357
func TestReferenceAttestation(t *testing.T) {
	sk := _byteArray("2c083f2c8fc923fa2bd32a70ab72b4b46247e8c1f347adc30b2f8036a355086c")
	pk := _byteArray("a9cf360aa15fb1d1d30ee2b578dc5884823c19661886ae8b892775ccb3bd96b7d7345569a2aa0b14e4d015c54a6a0c54")
	attestationDataByts := _byteArray("000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b0000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b")
	domain := _byteArray32("0100000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459")
	sig := _byteArray("b4fa352d2d6dbdf884266af7ea0914451929b343527ea6c1737ac93b3dde8b7c98e6ce61d68b7a2e7b7af8f8d0fd429d0bdd5f930b83e6842bf4342d3d1d3d10fc0d15bab7649bb8aa8287ca104a1f79d396ce0217bb5cd3e6503a3bce4c9776")

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(sk, "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, err)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	signer := NewSimpleSigner(wallet, &prot.NoProtection{}, core.PraterNetwork)

	// decode attestation
	attData := &phase0.AttestationData{}
	require.NoError(t, attData.UnmarshalSSZ(attestationDataByts))

	actualSig, root, err := signer.SignBeaconAttestation(attData, domain, pk)
	fmt.Println(string(root))
	require.NoError(t, err)
	require.EqualValues(t, sig, actualSig)
}

func TestAttestationSlashingSignatures(t *testing.T) {
	t.Run("valid attestation, sign using public key", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            67,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 0,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 1,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		}, _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})

	t.Run("valid attestation, sign using account name. Should error", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            67,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 0,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 1,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		}, _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			nil)
		require.NotNil(t, err)
		require.EqualError(t, err, "account was not supplied")
	})

	t.Run("double vote with different roots, should error", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		// first
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            67,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 0,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 1,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// second
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            67,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 0,
				Root:  _byteArray32("A"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 1,
				Root:  _byteArray32("A"),
			},
		},
			_byteArray32("A"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})

	t.Run("same vote with different domain, should not sign", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		// first
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            67,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 0,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 1,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// second
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            67,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 0,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 1,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01100000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})

	t.Run("surrounding vote, should err", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		// first
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            67,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 77,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 78,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// add another attestation building on the base
		// 77 <- 78 <- 79
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284116,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 78,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 79,
				Root:  _byteArray32("17959adc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// surround previous vote
		// 77 <- 78 <- 79
		// 	<- 80
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284117,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 77,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 80,
				Root:  _byteArray32("18959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})

	t.Run("surrounded vote, should err", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		// first
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284115,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 77,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 78,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// add another attestation building on the base
		// 77 <- 78 <----------------------100
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284116,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 78,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 100,
				Root:  _byteArray32("17959adc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// surround previous vote
		// 77 <- 78 <----------------------100
		// 								89 <- 90
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284117,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 89,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 90,
				Root:  _byteArray32("18959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})
}

func TestAttestationSignaturesNoSlashingData(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(t, seed, false, true)
	require.NoError(t, err)

	res, _, err := signer.SignBeaconAttestation(&phase0.AttestationData{
		Slot:            284115,
		Index:           2,
		BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
		Source: &phase0.Checkpoint{
			Epoch: 77,
			Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
		},
		Target: &phase0.Checkpoint{
			Epoch: 78,
			Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
		},
	},
		_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
		_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
	require.Nil(t, res)
	require.Error(t, err)
	require.EqualError(t, err, "highest attestation data is not found, can't determine if attestation is slashable")
}

func TestAttestationSignatures(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(t, seed, true, true)
	require.NoError(t, err)

	derivedSk, err := util.PrivateKeyFromSeedAndPath(seed, "m/12381/3600/0/0/0")
	require.NoError(t, err)

	sk := &bls.SecretKey{}
	require.NoError(t, sk.SetHexString(hex.EncodeToString(derivedSk.Marshal())))

	tests := []struct {
		name          string
		req           *phase0.AttestationData
		domain        [32]byte
		pubKey        []byte
		expectedError error
		accountPriv   *bls.SecretKey
		msg           string
	}{
		{
			name: "correct request",
			req: &phase0.AttestationData{
				Slot:            284115,
				Index:           2,
				BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
				Source: &phase0.Checkpoint{
					Epoch: 77,
					Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
				},
				Target: &phase0.Checkpoint{
					Epoch: 78,
					Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
				},
			},
			domain:        _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			expectedError: nil,
			accountPriv:   sk,
			msg:           "2783ca6dc161cc5feae0492ae79e52d7ae3eaff4b1f6b547d856533e9b733d8b",
		},
		{
			name: "far into the future source",
			req: &phase0.AttestationData{
				Slot:            284115,
				Index:           2,
				BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
				Source: &phase0.Checkpoint{
					Epoch: 1000077,
					Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
				},
				Target: &phase0.Checkpoint{
					Epoch: 78,
					Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
				},
			},
			domain:        _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			expectedError: errors.New("source epoch too far into the future"),
			accountPriv:   sk,
		},
		{
			name: "far into the future target",
			req: &phase0.AttestationData{
				Slot:            284115,
				Index:           2,
				BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
				Source: &phase0.Checkpoint{
					Epoch: 77,
					Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
				},
				Target: &phase0.Checkpoint{
					Epoch: 1000077,
					Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
				},
			},
			domain:        _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			expectedError: errors.New("target epoch too far into the future"),
			accountPriv:   sk,
		},
		{
			name: "unknown account, should error",
			req: &phase0.AttestationData{
				Slot:            284115,
				Index:           2,
				BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
				Source: &phase0.Checkpoint{
					Epoch: 77,
					Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
				},
				Target: &phase0.Checkpoint{
					Epoch: 78,
					Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
				},
			},
			domain:        _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			pubKey:        _byteArray("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3270"),
			expectedError: errors.New("account not found"),
			accountPriv:   nil,
			msg:           "",
		},
		{
			name: "nil account, should error",
			req: &phase0.AttestationData{
				Slot:            284115,
				Index:           2,
				BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
				Source: &phase0.Checkpoint{
					Epoch: 77,
					Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
				},
				Target: &phase0.Checkpoint{
					Epoch: 78,
					Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
				},
			},
			domain:        _byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			pubKey:        nil,
			expectedError: errors.New("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignBeaconAttestation(test.req, test.domain, test.pubKey)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)

				sig := bls.Sign{}
				require.NoError(t, sig.Deserialize(res))

				msgBytes, err := hex.DecodeString(test.msg)
				require.NoError(t, err)
				require.True(t, sig.VerifyByte(test.accountPriv.GetPublicKey(), msgBytes))
			}
		})
	}
}

func TestFarFutureAttestationSignature(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	network := core.PraterNetwork
	maxValidEpoch := network.EstimatedEpochAtSlot(network.EstimatedSlotAtTime(time.Now().Unix() + FarFutureMaxValidEpoch))

	t.Run("max valid source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284115,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: maxValidEpoch,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 78,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})
	t.Run("too far into the future source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284115,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: maxValidEpoch + 1,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 78,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.EqualError(t, err, "source epoch too far into the future")
	})
	t.Run("max valid target", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)
		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284115,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 77,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: maxValidEpoch,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))

		require.NoError(t, err)
	})
	t.Run("too far into the future target", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		_, _, err = signer.SignBeaconAttestation(&phase0.AttestationData{
			Slot:            284115,
			Index:           2,
			BeaconBlockRoot: _byteArray32("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e"),
			Source: &phase0.Checkpoint{
				Epoch: 77,
				Root:  _byteArray32("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d"),
			},
			Target: &phase0.Checkpoint{
				Epoch: maxValidEpoch + 1,
				Root:  _byteArray32("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0"),
			},
		},
			_byteArray32("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac"),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))

		require.EqualError(t, err, "target epoch too far into the future")
	})
}
