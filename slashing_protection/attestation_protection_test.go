package slashing_protection

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

func setupAttestation(t *testing.T) (core.SlashingProtector, []core.ValidatorAccount) {
	err := e2types.InitBLS()
	require.NoError(t, err)

	// seed
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	// create an account to use
	vault, err := vault()
	require.NoError(t, err)

	w, err := vault.Wallet()
	require.NoError(t, err)

	account1, err := w.CreateValidatorAccount(seed, nil)
	require.NoError(t, err)

	account2, err := w.CreateValidatorAccount(seed, nil)
	require.NoError(t, err)

	protector := NewNormalProtection(vault.Context.Storage.(core.SlashingStore))
	err = protector.SaveAttestation(account1.ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.AttestationData{
			Slot:            30,
			CommitteeIndex:  5,
			BeaconBlockRoot: []byte("A"),
			Source: &pb.Checkpoint{
				Epoch: 1,
				Root:  []byte("B"),
			},
			Target: &pb.Checkpoint{
				Epoch: 2,
				Root:  []byte("C"),
			},
		},
	})
	require.NoError(t, err)

	err = protector.SaveAttestation(account1.ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.AttestationData{
			Slot:            30,
			CommitteeIndex:  5,
			BeaconBlockRoot: []byte("A"),
			Source: &pb.Checkpoint{
				Epoch: 2,
				Root:  []byte("B"),
			},
			Target: &pb.Checkpoint{
				Epoch: 3,
				Root:  []byte("C"),
			},
		},
	})
	require.NoError(t, err)

	err = protector.SaveAttestation(account1.ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.AttestationData{
			Slot:            30,
			CommitteeIndex:  5,
			BeaconBlockRoot: []byte("B"),
			Source: &pb.Checkpoint{
				Epoch: 3,
				Root:  []byte("C"),
			},
			Target: &pb.Checkpoint{
				Epoch: 4,
				Root:  []byte("D"),
			},
		},
	})
	require.NoError(t, err)

	err = protector.SaveAttestation(account1.ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.AttestationData{
			Slot:            30,
			CommitteeIndex:  5,
			BeaconBlockRoot: []byte("B"),
			Source: &pb.Checkpoint{
				Epoch: 4,
				Root:  []byte("C"),
			},
			Target: &pb.Checkpoint{
				Epoch: 10,
				Root:  []byte("D"),
			},
		},
	})
	require.NoError(t, err)

	err = protector.SaveAttestation(account1.ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.AttestationData{
			Slot:            30,
			CommitteeIndex:  5,
			BeaconBlockRoot: []byte("B"),
			Source: &pb.Checkpoint{
				Epoch: 5,
				Root:  []byte("C"),
			},
			Target: &pb.Checkpoint{
				Epoch: 9,
				Root:  []byte("D"),
			},
		},
	})
	require.NoError(t, err)

	return protector, []core.ValidatorAccount{account1, account2}
}

func TestSurroundingVote(t *testing.T) {
	protector, accounts := setupAttestation(t)

	t.Run("1 Surrounded vote", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  4,
				BeaconBlockRoot: []byte("A"),
				Source: &pb.Checkpoint{
					Epoch: 2,
					Root:  []byte("B"),
				},
				Target: &pb.Checkpoint{
					Epoch: 5,
					Root:  []byte("C"),
				},
			},
		})

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, core.SurroundingVote, res[0].Status)
	})

	t.Run("2 Surrounded votes", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  4,
				BeaconBlockRoot: []byte("A"),
				Source: &pb.Checkpoint{
					Epoch: 1,
					Root:  []byte("B"),
				},
				Target: &pb.Checkpoint{
					Epoch: 7,
					Root:  []byte("C"),
				},
			},
		})

		require.NoError(t, err)
		require.Len(t, res, 2)
		require.False(t, res[0].Status != core.SurroundingVote || res[1].Status != core.SurroundingVote)
	})

	t.Run("1 Surrounding vote", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  4,
				BeaconBlockRoot: []byte("A"),
				Source: &pb.Checkpoint{
					Epoch: 5,
					Root:  []byte("B"),
				},
				Target: &pb.Checkpoint{
					Epoch: 7,
					Root:  []byte("C"),
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, core.SurroundedVote, res[0].Status)
	})

	t.Run("2 Surrounding vote", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  4,
				BeaconBlockRoot: []byte("A"),
				Source: &pb.Checkpoint{
					Epoch: 6,
					Root:  []byte("B"),
				},
				Target: &pb.Checkpoint{
					Epoch: 7,
					Root:  []byte("C"),
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, res, 2)
		require.False(t, res[0].Status != core.SurroundedVote || res[1].Status != core.SurroundedVote)
	})
}

func TestDoubleAttestationVote(t *testing.T) {
	protector, accounts := setupAttestation(t)

	t.Run("Different committee index, should slash", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  4,
				BeaconBlockRoot: []byte("A"),
				Source: &pb.Checkpoint{
					Epoch: 2,
					Root:  []byte("B"),
				},
				Target: &pb.Checkpoint{
					Epoch: 3,
					Root:  []byte("C"),
				},
			},
		})

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, core.DoubleVote, res[0].Status)
	})

	t.Run("Different block root, should slash", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  5,
				BeaconBlockRoot: []byte("AA"),
				Source: &pb.Checkpoint{
					Epoch: 2,
					Root:  []byte("B"),
				},
				Target: &pb.Checkpoint{
					Epoch: 3,
					Root:  []byte("C"),
				},
			},
		})

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, core.DoubleVote, res[0].Status)
	})

	t.Run("Same attestation, should not error", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  5,
				BeaconBlockRoot: []byte("B"),
				Source: &pb.Checkpoint{
					Epoch: 3,
					Root:  []byte("C"),
				},
				Target: &pb.Checkpoint{
					Epoch: 4,
					Root:  []byte("D"),
				},
			},
		})
		require.False(t, err != nil || len(res) != 0)
	})

	t.Run("new attestation, should not error", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.AttestationData{
				Slot:            30,
				CommitteeIndex:  5,
				BeaconBlockRoot: []byte("E"),
				Source: &pb.Checkpoint{
					Epoch: 10,
					Root:  []byte("I"),
				},
				Target: &pb.Checkpoint{
					Epoch: 11,
					Root:  []byte("H"),
				},
			},
		})
		require.False(t, err != nil || len(res) != 0)
	})
}
