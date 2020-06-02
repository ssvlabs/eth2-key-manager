package validator_signer

import (
	"encoding/hex"
	"fmt"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	"testing"
)

func TestAttestationSlashingSignatures(t *testing.T) {
	seed,_ := hex.DecodeString("f51883a4c56467458c3b47d06cd135f862a6266fabdfb9e9e4702ea5511375d7")
	signer,err := setupWithSlashingProtection(seed)
	if err != nil {
		t.Error(err)
		return
	}

	_,err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
		Id:                   &pb.SignBeaconAttestationRequest_Account{Account:"1"},
		Domain:               ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
		Data:                 &pb.AttestationData{
			Slot:            284115,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
			Source:          &pb.Checkpoint{
				Epoch:                8877,
				Root:                 ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
			},
			Target:          &pb.Checkpoint{
				Epoch:                8878,
				Root:                 ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	// double vote with different root
	_,err = signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
		Id:                   &pb.SignBeaconAttestationRequest_Account{Account:"1"},
		Domain:               ignoreError(hex.DecodeString("A")).([]byte),
		Data:                 &pb.AttestationData{
			Slot:            284115,
			CommitteeIndex:  2,
			BeaconBlockRoot: ignoreError(hex.DecodeString("A")).([]byte),
			Source:          &pb.Checkpoint{
				Epoch:                8877,
				Root:                 ignoreError(hex.DecodeString("A")).([]byte),
			},
			Target:          &pb.Checkpoint{
				Epoch:                8878,
				Root:                 ignoreError(hex.DecodeString("A")).([]byte),
			},
		},
	})
	expectedErr := "slashable attestation, not signing"
	if err == nil {
		t.Errorf("expectd an error, did not error")
	} else if err.Error() != expectedErr {
		t.Errorf("received error: %s, different than expected: %s", err.Error(), expectedErr)
	}
}

func TestAttestationSignatures(t *testing.T) {
	seed,_ := hex.DecodeString("f51883a4c56467458c3b47d06cd135f862a6266fabdfb9e9e4702ea5511375d7")
	signer,err := setupWithSlashingProtection(seed)
	if err != nil {
		t.Error(err)
		return
	}
	accountPriv,err := util.PrivateKeyFromSeedAndPath(seed,"m/12381/3600/0/0")
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string
		req *pb.SignBeaconAttestationRequest
		expectedError error
		accountPriv *e2types.BLSPrivateKey
		msg string
	}{
		{
			name:"correct request",
			req: &pb.SignBeaconAttestationRequest{
				Id:                   &pb.SignBeaconAttestationRequest_Account{Account:"1"},
				Domain:               ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data:                 &pb.AttestationData{
					Slot:            284115,
					CommitteeIndex:  2,
					BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
					Source:          &pb.Checkpoint{
						Epoch:                8877,
						Root:                 ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
					},
					Target:          &pb.Checkpoint{
						Epoch:                8878,
						Root:                 ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
					},
				},
			},
			expectedError:nil,
			accountPriv:accountPriv,
			msg:"6c66b61134300a3eeb37b0788bd8fc32663e3ada6b8d2e1fc7801641a3851300",
		},
		{
			name:"unknown account, should error",
			req: &pb.SignBeaconAttestationRequest{
				Id:                   &pb.SignBeaconAttestationRequest_Account{Account:"10"},
				Domain:               ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data:                 &pb.AttestationData{
					Slot:            284115,
					CommitteeIndex:  2,
					BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
					Source:          &pb.Checkpoint{
						Epoch:                8877,
						Root:                 ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
					},
					Target:          &pb.Checkpoint{
						Epoch:                8878,
						Root:                 ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
					},
				},
			},
			expectedError:fmt.Errorf("no account with name \"10\""),
			accountPriv:nil,
			msg:"",
		},
		{
			name:"unable to unlock account, should error",
			req: &pb.SignBeaconAttestationRequest{
				Id:                   &pb.SignBeaconAttestationRequest_Account{Account:"2"},
				Domain:               ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data:                 &pb.AttestationData{
					Slot:            284115,
					CommitteeIndex:  2,
					BeaconBlockRoot: ignoreError(hex.DecodeString("7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e")).([]byte),
					Source:          &pb.Checkpoint{
						Epoch:                8877,
						Root:                 ignoreError(hex.DecodeString("7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d")).([]byte),
					},
					Target:          &pb.Checkpoint{
						Epoch:                8878,
						Root:                 ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
					},
				},
			},
			expectedError:fmt.Errorf("incorrect passphrase"),
			accountPriv:nil,
			msg:"",
		},
	}

	for _,test := range tests {
		t.Run(test.name,func(t *testing.T) {
			res,err := signer.SignBeaconAttestation(test.req)
			if test.expectedError != nil {
				if err != nil {
					if err.Error() != test.expectedError.Error() {
						t.Errorf("wrong error returned: %s, expected: %s", err.Error(),test.expectedError.Error())
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

				sig,err := e2types.BLSSignatureFromBytes(res.Signature)
				if err != nil {
					t.Error(err)
					return
				}
				msgBytes,err := hex.DecodeString(test.msg)
				if err != nil {
					t.Error(err)
					return
				}
				if !sig.Verify(msgBytes,test.accountPriv.PublicKey()) {
					t.Errorf("signature does not verify against pubkey",)
				}
			}
		})
	}
}
