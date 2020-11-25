package validator_signer

import (
	"encoding/hex"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/prysmaticlabs/prysm/shared/timeutils"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func TestProposalSlashingSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(seed, true)
	require.NoError(t, err)

	t.Run("valid proposal", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NoError(t, err)
	})

	t.Run("valid proposal, sign using account name. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_Account{Account: "1"},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.Error(t, err, "account was not supplied")
	})

	t.Run("double proposal, different state root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("A"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different body root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("A"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different parent root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different proposer index. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 3,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})
}

func TestFarFutureProposalSignature(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	network := core.PyrmontNetwork
	maxValidSlot := network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch)

	t.Run("max valid source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          maxValidSlot,
				ProposerIndex: 3,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NoError(t, err)
	})
	t.Run("too far into the future source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)
		_, err = signer.SignBeaconProposal(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          maxValidSlot+1,
				ProposerIndex: 3,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.EqualError(t, err, "proposed block slot too far into the future")
	})
}