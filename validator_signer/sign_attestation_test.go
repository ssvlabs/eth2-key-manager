package validator_signer

import (
	"encoding/hex"
	"testing"

	"github.com/bloxapp/eth2-key-manager/core"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/timeutils"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func ignoreError(val interface{}, err error) interface{} {
	return val
}

func TestAttestationSlashingSignatures(t *testing.T) {
	t.Run("valid attestation, sign using public key", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)

		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            67,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 1,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		}, ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})

	t.Run("valid attestation, sign using account name. Should error", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            67,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 1,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		}, ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			nil)
		require.NotNil(t, err)
		require.EqualError(t, err, "account was not supplied")
	})

	t.Run("double vote with different roots, should error", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            67,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 1,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// second
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            67,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("A")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 1,
				Root:  ignoreError(hex.DecodeString("A")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("A")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})

	t.Run("same vote with different domain, should not sign", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            67,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 1,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// second
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            67,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 1,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01100000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})

	t.Run("surrounding vote, should err", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            67,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 77,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 78,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// add another attestation building on the base
		// 77 <- 78 <- 79
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284116,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 78,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 79,
				Root:  ignoreError(hex.DecodeString("17959adc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// surround previous vote
		// 77 <- 78 <- 79
		// 	<- 80
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284117,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 77,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 80,
				Root:  ignoreError(hex.DecodeString("18959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})

	t.Run("surrounded vote, should err", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284115,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 77,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 78,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// add another attestation building on the base
		// 77 <- 78 <----------------------100
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284116,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 78,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 100,
				Root:  ignoreError(hex.DecodeString("17959adc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)

		// surround previous vote
		// 77 <- 78 <----------------------100
		// 								89 <- 90
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284117,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 89,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 90,
				Root:  ignoreError(hex.DecodeString("18959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
	})
}

func TestAttestationSignaturesNoSlashingData(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(seed, false)
	require.NoError(t, err)

	res, err := signer.SignBeaconAttestation(&eth.AttestationData{
		Slot:            284115,
		CommitteeIndex:  2,
		BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
		Source: &eth.Checkpoint{
			Epoch: 77,
			Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
		},
		Target: &eth.Checkpoint{
			Epoch: 78,
			Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
		},
	},
		ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
		_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
	require.Nil(t, res)
	require.EqualError(t, err, "highest attestation data is nil, can't determine if attestation is slashable")
}

func TestAttestationSignatures(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(seed, true)
	require.NoError(t, err)

	accountPriv, err := util.PrivateKeyFromSeedAndPath(seed, "m/12381/3600/0/0/0")
	require.NoError(t, err)

	tests := []struct {
		name          string
		req           *eth.AttestationData
		domain        []byte
		pubKey        []byte
		expectedError error
		accountPriv   *e2types.BLSPrivateKey
		msg           string
	}{
		{
			name: "correct request",
			req: &eth.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &eth.Checkpoint{
					Epoch: 77,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &eth.Checkpoint{
					Epoch: 78,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
			domain:        ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			expectedError: nil,
			accountPriv:   accountPriv,
			msg:           "2783ca6dc161cc5feae0492ae79e52d7ae3eaff4b1f6b547d856533e9b733d8b",
		},
		{
			name: "far into the future source",
			req: &eth.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &eth.Checkpoint{
					Epoch: 1000077,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &eth.Checkpoint{
					Epoch: 78,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
			domain:        ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			expectedError: errors.New("source epoch too far into the future"),
			accountPriv:   accountPriv,
		},
		{
			name: "far into the future target",
			req: &eth.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &eth.Checkpoint{
					Epoch: 77,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &eth.Checkpoint{
					Epoch: 1000077,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
			domain:        ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			expectedError: errors.New("target epoch too far into the future"),
			accountPriv:   accountPriv,
		},
		{
			name: "unknown account, should error",
			req: &eth.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &eth.Checkpoint{
					Epoch: 77,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &eth.Checkpoint{
					Epoch: 78,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
			domain:        ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			pubKey:        _byteArray("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3270"),
			expectedError: errors.New("account not found"),
			accountPriv:   nil,
			msg:           "",
		},
		{
			name: "nil account, should error",
			req: &eth.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &eth.Checkpoint{
					Epoch: 77,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &eth.Checkpoint{
					Epoch: 78,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
			domain:        ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			pubKey:        nil,
			expectedError: errors.New("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.SignBeaconAttestation(test.req, test.domain, test.pubKey)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)

				sig, err := e2types.BLSSignatureFromBytes(res)
				require.NoError(t, err)

				msgBytes, err := hex.DecodeString(test.msg)
				require.NoError(t, err)
				require.True(t, sig.Verify(msgBytes, test.accountPriv.PublicKey()))
			}
		})
	}
}

func TestFarFutureAttestationSignature(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	network := core.PyrmontNetwork
	maxValidEpoch := network.EstimatedEpochAtSlot(network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch))

	t.Run("max valid source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)

		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284115,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: maxValidEpoch,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 78,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})
	t.Run("too far into the future source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284115,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: maxValidEpoch + 1,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 78,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.EqualError(t, err, "source epoch too far into the future")
	})
	t.Run("max valid target", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)
		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284115,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 77,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: maxValidEpoch,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))

		require.NoError(t, err)
	})
	t.Run("too far into the future target", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)

		_, err = signer.SignBeaconAttestation(&eth.AttestationData{
			Slot:            284115,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 77,
				Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: maxValidEpoch + 1,
				Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
			ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))

		require.EqualError(t, err, "target epoch too far into the future")
	})
}
