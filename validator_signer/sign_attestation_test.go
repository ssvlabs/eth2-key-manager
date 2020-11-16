package validator_signer

import (
	"encoding/hex"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func TestAttestationSlashingSignatures(t *testing.T) {
	t.Run("valid attestation, sign using public key", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed)
		require.NoError(t, err)
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("valid attestation, sign using account name. Should error", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed)
		require.NoError(t, err)
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_Account{Account: "1"},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "account was not supplied")
	})

	t.Run("double vote with different roots, should error", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)

		// second
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("A")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("A")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("A")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("A")).([]byte),
				},
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (DoubleVote), not signing")
	})

	t.Run("same vote with different domain, should sign", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)

		// second
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01100000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("surrounding vote, should err", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)

		// add another attestation building on the base
		// 8877 <- 8878 <- 8879
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284116,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8879,
					Root:  ignoreError(hex.DecodeString("17959adc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)

		// surround previous vote
		// 8877 <- 8878 <- 8879
		// 	<- 8880
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284117,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8880,
					Root:  ignoreError(hex.DecodeString("18959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (SurroundingVote), not signing")
	})

	t.Run("surrounded vote, should err", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed)
		require.NoError(t, err)

		// first
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284115,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8877,
					Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)

		// add another attestation building on the base
		// 8877 <- 8878 <----------------------9000
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284116,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8878,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 9000,
					Root:  ignoreError(hex.DecodeString("17959adc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NoError(t, err)

		// surround previous vote
		// 8877 <- 8878 <- 8879 <----------------------9000
		// 								8900 <- 8901
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data: &pb.AttestationData{
				Slot:            284117,
				CommitteeIndex:  2,
				BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
				Source: &pb.Checkpoint{
					Epoch: 8900,
					Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
				Target: &pb.Checkpoint{
					Epoch: 8901,
					Root:  ignoreError(hex.DecodeString("18959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable attestation (SurroundedVote), not signing")
	})
}

func TestAttestationSignatures(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(seed)
	require.NoError(t, err)

	accountPriv, err := util.PrivateKeyFromSeedAndPath(seed, "m/12381/3600/0/0/0")
	require.NoError(t, err)

	tests := []struct {
		name          string
		req           *pb.SignBeaconAttestationRequest
		expectedError error
		accountPriv   *e2types.BLSPrivateKey
		msg           string
	}{
		{
			name: "correct request",
			req: &pb.SignBeaconAttestationRequest{
				Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
				Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data: &pb.AttestationData{
					Slot:            284115,
					CommitteeIndex:  2,
					BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
					Source: &pb.Checkpoint{
						Epoch: 8877,
						Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
					},
					Target: &pb.Checkpoint{
						Epoch: 8878,
						Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
					},
				},
			},
			expectedError: nil,
			accountPriv:   accountPriv,
			msg:           "6c66b61134300a3eeb37b0788bd8fc32663e3ada6b8d2e1fc7801641a3851300",
		},
		{
			name: "unknown account, should error",
			req: &pb.SignBeaconAttestationRequest{
				Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3270")},
				Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data: &pb.AttestationData{
					Slot:            284115,
					CommitteeIndex:  2,
					BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
					Source: &pb.Checkpoint{
						Epoch: 8877,
						Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
					},
					Target: &pb.Checkpoint{
						Epoch: 8878,
						Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
					},
				},
			},
			expectedError: errors.New("account not found"),
			accountPriv:   nil,
			msg:           "",
		},
		{
			name: "nil account, should error",
			req: &pb.SignBeaconAttestationRequest{
				Id:     nil,
				Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data: &pb.AttestationData{
					Slot:            284115,
					CommitteeIndex:  2,
					BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
					Source: &pb.Checkpoint{
						Epoch: 8877,
						Root:  ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
					},
					Target: &pb.Checkpoint{
						Epoch: 8878,
						Root:  ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
					},
				},
			},
			expectedError: errors.New("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.SignBeaconAttestation(test.req)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)

				sig, err := e2types.BLSSignatureFromBytes(res.Signature)
				require.NoError(t, err)

				msgBytes, err := hex.DecodeString(test.msg)
				require.NoError(t, err)
				require.True(t, sig.Verify(msgBytes, test.accountPriv.PublicKey()))
			}
		})
	}
}
