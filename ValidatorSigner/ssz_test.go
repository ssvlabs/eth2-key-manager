package ValidatorSigner

import (
	"encoding/hex"
	"fmt"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	"testing"
)

func ignoreError(val interface{}, err error)interface{} {
	return val
}

func proposalFixtures() ([]*pb.SignBeaconProposalRequest,[]string) {
	return []*pb.SignBeaconProposalRequest{
			{
				Id:                   nil,
				Domain:               ignoreError(hex.DecodeString("00000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data:                 &pb.BeaconBlockHeader{
					Slot:                 285193,
					ProposerIndex:        15091,
					ParentRoot:           ignoreError(hex.DecodeString("ccc2dc43cc3b39a7ab994fd1cb269508698580fd4c4b145a2a75873c331fce7e")).([]byte),
					StateRoot:            ignoreError(hex.DecodeString("0c683084d760a21305926b83de104d75319033658efb7c1d5067fb5ff97b8794")).([]byte),
					BodyRoot:             ignoreError(hex.DecodeString("c0044ddb3f90d362640ce1fa68a07caee31814008b6d7e0eb48dd56e1a77e13e")).([]byte),
				},
			},
			{
				Id:                   nil,
				Domain:               ignoreError(hex.DecodeString("00000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data:                 &pb.BeaconBlockHeader{
					Slot:                 285199,
					ProposerIndex:        35214,
					ParentRoot:           ignoreError(hex.DecodeString("c47bee4e45beebee094e675a6109e7a59db58ccc4548fd884be76868314adffe")).([]byte),
					StateRoot:            ignoreError(hex.DecodeString("32d52152130f9759e420779b26ae91ed02ccc81b3e0997438d411a1fd44dd9aa")).([]byte),
					BodyRoot:             ignoreError(hex.DecodeString("7fa4feb590dbad29f691234890f98f26b3d38a451cb9b991443ca4fc5fcafc0b")).([]byte),
				},
			},
			{
				Id:                   nil,
				Domain:               ignoreError(hex.DecodeString("00000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
				Data:                 &pb.BeaconBlockHeader{
					Slot:                 285200,
					ProposerIndex:        35016,
					ParentRoot:           ignoreError(hex.DecodeString("04b35797226b8c338e5dd0dbe47f0d60b9b5248b2bdc70691b04a3772a188628")).([]byte),
					StateRoot:            ignoreError(hex.DecodeString("cc0fd17797417ef263be5dd94eccc0040047ce122507215f69ca6504e3ebcc6e")).([]byte),
					BodyRoot:             ignoreError(hex.DecodeString("b861ef8113ff7221dbf4d1d2e2cd24ed9c1a1d1aae2819b22720e090a5904eef")).([]byte),
				},
			},
	},
	[]string{
		"9f68c38293403d2e272749ab0089d4ef48ca535695dc1a3a1991fca4ecb8d223",
		"995020baa0c37625c9f64cf5f929f0f8f30c544e37bfab6e9d07c54c284a5ab7",
		"957a2c315b2e47c8333c1545b269b874cdcf2758525cf1f87d708c42162a1ca7",
	}
}

func attestationsFixtures() ([]*pb.SignBeaconAttestationRequest,[]string) {
	return []*pb.SignBeaconAttestationRequest{
		{
			Id:                   nil,
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
		{
			Id:                   nil,
			Domain:               ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data:                 &pb.AttestationData{
				Slot:            284147,
				CommitteeIndex:  0,
				BeaconBlockRoot: ignoreError(hex.DecodeString("de68bd3626a43349f6134acaccd721e9858381051bc8488a6f9b1d3a8bcc75b2")).([]byte),
				Source:          &pb.Checkpoint{
					Epoch:                8878,
					Root:                 ignoreError(hex.DecodeString("17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0")).([]byte),
				},
				Target:          &pb.Checkpoint{
					Epoch:                8879,
					Root:                 ignoreError(hex.DecodeString("0fa1c8f5c8008b97bc36fa54b3c17325decdd22829e5bc2c2939f095776c6ce6")).([]byte),
				},
			},
		},
		{
			Id:                   nil,
			Domain:               ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
			Data:                 &pb.AttestationData{
				Slot:            284178,
				CommitteeIndex:  6,
				BeaconBlockRoot: ignoreError(hex.DecodeString("10720f9ce0e41c7b099b32be41a5ac62f3abdc25bf1f15fc5f0e7a6eb258301f")).([]byte),
				Source:          &pb.Checkpoint{
					Epoch:                8879,
					Root:                 ignoreError(hex.DecodeString("0fa1c8f5c8008b97bc36fa54b3c17325decdd22829e5bc2c2939f095776c6ce6")).([]byte),
				},
				Target:          &pb.Checkpoint{
					Epoch:                8880,
					Root:                 ignoreError(hex.DecodeString("f52eaa102ec0d310e0cb580261398d24e5295ad5740c8614af1cbb86ec248b3c")).([]byte),
				},
			},
		},

	}, []string{
			"6c66b61134300a3eeb37b0788bd8fc32663e3ada6b8d2e1fc7801641a3851300",
			"a69d9159fb2be639ef74749eca81d60490e1589bf75011aca8c5c0324cf463c0",
			"310351fd11e8fc3c7342ca91fde5e29b69873df4f378355988af583cd47b9206",
		}
}

func TestAttestationRootComputation(t *testing.T) {
	attestations, roots := attestationsFixtures()

	for i := range attestations {
		t.Run(fmt.Sprintf("TestRootComputation %d", i), func(t *testing.T) {
			root,error := prepareAttestationReqForSigning(attestations[i])
			if error != nil {
				t.Error(error)
			}
			expectedroot := roots[i]
			if val := hex.EncodeToString(root); val != expectedroot {
				t.Error(fmt.Errorf("calculateed root: %s not equal to expected root: %s",val, expectedroot))
			}
		})
	}
}


func TestProposalRootComputation(t *testing.T) {
	proposals, roots := proposalFixtures()
	for i := range proposals {
		t.Run(fmt.Sprintf("TestRootComputation %d", i), func(t *testing.T) {
			root,error := prepareProposalReqForSigning(proposals[i])
			if error != nil {
				t.Error(error)
			}
			expectedroot := roots[i]
			if val := hex.EncodeToString(root); val != expectedroot {
				t.Error(fmt.Errorf("calculateed root: %s not equal to expected root: %s",val, expectedroot))
			}
		})
	}
}

