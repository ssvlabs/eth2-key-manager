package slashing_protection

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func store() *in_memory.InMemStore {
	return in_memory.NewInMemStore(core.MainNetwork)
}

func vault() (*eth2keymanager.KeyVault, error) {
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store())
	return eth2keymanager.NewKeyVault(options)
}

func setupProposal() (core.SlashingProtector, []core.ValidatorAccount, error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil, nil, err
	}

	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	// create an account to use
	vault, err := vault()
	if err != nil {
		return nil, nil, err
	}
	w, err := vault.Wallet()
	if err != nil {
		return nil, nil, err
	}
	account1, err := w.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, nil, err
	}
	account2, err := w.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, nil, err
	}

	protector := NewNormalProtection(vault.Context.Storage.(core.SlashingStore))
	protector.SaveProposal(account1.ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.BeaconBlockHeader{
			Slot:          100,
			ProposerIndex: 2,
			ParentRoot:    []byte("A"),
			StateRoot:     []byte("A"),
			BodyRoot:      []byte("A"),
		},
	})
	protector.SaveProposal(account1.ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.BeaconBlockHeader{
			Slot:          101,
			ProposerIndex: 2,
			ParentRoot:    []byte("B"),
			StateRoot:     []byte("B"),
			BodyRoot:      []byte("B"),
		},
	})
	protector.SaveProposal(account1.ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
		Id:     nil,
		Domain: []byte("domain"),
		Data: &pb.BeaconBlockHeader{
			Slot:          102,
			ProposerIndex: 2,
			ParentRoot:    []byte("C"),
			StateRoot:     []byte("C"),
			BodyRoot:      []byte("C"),
		},
	})

	return protector, []core.ValidatorAccount{account1, account2}, nil
}

func TestDoubleProposal(t *testing.T) {
	protector, accounts, err := setupProposal()
	require.NoError(t, err)

	t.Run("New proposal, should not slash", func(t *testing.T) {
		res := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.Equal(t, res.Status, core.ValidProposal)
	})

	t.Run("different proposer index, should not slash", func(t *testing.T) {
		res := protector.IsSlashableProposal(accounts[1].ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.BeaconBlockHeader{
				Slot:          100,
				ProposerIndex: 3,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("A"),
				BodyRoot:      []byte("A"),
			},
		})
		require.Equal(t, res.Status, core.ValidProposal)
	})

	t.Run("double proposal (different body root), should slash", func(t *testing.T) {
		res := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.BeaconBlockHeader{
				Slot:          100,
				ProposerIndex: 2,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("A"),
				BodyRoot:      []byte("B"),
			},
		})
		require.Equal(t, res.Status, core.DoubleProposal)
	})

	t.Run("double proposal (different state root), should slash", func(t *testing.T) {
		res := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.BeaconBlockHeader{
				Slot:          100,
				ProposerIndex: 2,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("B"),
				BodyRoot:      []byte("A"),
			},
		})
		require.Equal(t, res.Status, core.DoubleProposal)
	})

	t.Run("double proposal (different state and body root), should slash", func(t *testing.T) {
		res := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), &pb.SignBeaconProposalRequest{
			Id:     nil,
			Domain: []byte("domain"),
			Data: &pb.BeaconBlockHeader{
				Slot:          100,
				ProposerIndex: 2,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("B"),
				BodyRoot:      []byte("B"),
			},
		})
		require.Equal(t, res.Status, core.DoubleProposal)
	})
}
