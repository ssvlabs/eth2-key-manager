package validator_signer

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func TestAttestationSlashingSignatures(t *testing.T) {
	t.Run("valid attestation, sign using account name", func(t *testing.T) {
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
		require.NoError(t, err)
	})

	t.Run("valid attestation, sign using pub key. Should error", func(t *testing.T) {
		seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
		signer, err := setupWithSlashingProtection(seed)
		require.NoError(t, err)
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("a279033cc76667b4d083a605b7656ee48629c9e22032fb2a631b8e2c025c7000b87fc9fa5df47e107b51f436749d38ab")},
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
		signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
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

		// second
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_Account{Account: "1"},
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
		signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
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

		// second
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_Account{Account: "1"},
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
		signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
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


		// add another attestation building on the base
		// 8877 <- 8878 <- 8879
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_Account{Account: "1"},
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
			Id:     &pb.SignBeaconAttestationRequest_Account{Account: "1"},
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
		signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("a279033cc76667b4d083a605b7656ee48629c9e22032fb2a631b8e2c025c7000b87fc9fa5df47e107b51f436749d38ab")},
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

		// add another attestation building on the base
		// 8877 <- 8878 <----------------------9000
		_, err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
			Id:     &pb.SignBeaconAttestationRequest_Account{Account: "1"},
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
			Id:     &pb.SignBeaconAttestationRequest_Account{Account: "1"},
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
	if err != nil {
		t.Error(err)
		return
	}
	accountPriv, err := util.PrivateKeyFromSeedAndPath(seed, "m/12381/3600/0/0/0")
	if err != nil {
		t.Error(err)
		return
	}

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
			},
			expectedError: nil,
			accountPriv:   accountPriv,
			msg:           "6c66b61134300a3eeb37b0788bd8fc32663e3ada6b8d2e1fc7801641a3851300",
		},
		{
			name: "unknown account, should error",
			req: &pb.SignBeaconAttestationRequest{
				Id:     &pb.SignBeaconAttestationRequest_Account{Account: "10"},
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
			expectedError: fmt.Errorf("account not found"),
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
			expectedError: fmt.Errorf("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.SignBeaconAttestation(test.req)
			if test.expectedError != nil {
				if err != nil {
					if err.Error() != test.expectedError.Error() {
						t.Errorf("wrong error returned: %s, expected: %s", err.Error(), test.expectedError.Error())
					}
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				if err != nil {
					t.Error(err)
					return
				}

				sig, err := e2types.BLSSignatureFromBytes(res.Signature)
				if err != nil {
					t.Error(err)
					return
				}
				msgBytes, err := hex.DecodeString(test.msg)
				if err != nil {
					t.Error(err)
					return
				}
				if !sig.Verify(msgBytes, test.accountPriv.PublicKey()) {
					t.Errorf("signature does not verify against pubkey", )
				}
			}
		})
	}
}
