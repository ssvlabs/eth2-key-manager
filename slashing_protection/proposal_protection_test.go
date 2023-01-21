package slashingprotection

import (
	"encoding/hex"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
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

func store() *inmemory.InMemStore {
	return inmemory.NewInMemStore(core.MainNetwork)
}

func vault() (*eth2keymanager.KeyVault, error) {
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store())
	return eth2keymanager.NewKeyVault(options)
}

func setupProposal(t *testing.T, updateHighestProposal bool) (core.SlashingProtector, []core.ValidatorAccount, error) {
	require.NoError(t, core.InitBLS()) // very important!!!

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

	if updateHighestProposal {
		blk := &phase0.BeaconBlock{
			Slot:          100,
			ProposerIndex: 2,
			ParentRoot:    _byteArray32("A"),
			StateRoot:     _byteArray32("A"),
			Body:          &phase0.BeaconBlockBody{},
		}
		require.NoError(t, protector.UpdateHighestProposal(account1.ValidatorPublicKey(), blk.Slot))
	}

	return protector, []core.ValidatorAccount{account1, account2}, nil
}

func TestProposalProtection(t *testing.T) {
	t.Run("New proposal, should not slash", func(t *testing.T) {
		protector, accounts, err := setupProposal(t, true)
		require.NoError(t, err)

		blk := &phase0.BeaconBlock{
			Slot:          101,
			ProposerIndex: 2,
			ParentRoot:    _byteArray32("Z"),
			StateRoot:     _byteArray32("Z"),
			Body:          &phase0.BeaconBlockBody{},
		}

		res, err := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), blk.Slot)
		require.NoError(t, err)
		require.Equal(t, res.Status, core.ValidProposal)
	})

	t.Run("No highest proposal db, should error", func(t *testing.T) {
		protector, accounts, err := setupProposal(t, false)
		require.NoError(t, err)

		blk := &phase0.BeaconBlock{
			Slot:          99,
			ProposerIndex: 2,
			ParentRoot:    _byteArray32("Z"),
			StateRoot:     _byteArray32("Z"),
			Body:          &phase0.BeaconBlockBody{},
		}
		res, err := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), blk.Slot)
		require.EqualError(t, err, "highest proposal data is nil, can't determine if proposal is slashable")
		require.Nil(t, res)
	})

	t.Run("Lower than highest proposal db, should error", func(t *testing.T) {
		protector, accounts, err := setupProposal(t, true)
		require.NoError(t, err)

		blk := &phase0.BeaconBlock{
			Slot:          99,
			ProposerIndex: 2,
			ParentRoot:    _byteArray32("Z"),
			StateRoot:     _byteArray32("Z"),
			Body:          &phase0.BeaconBlockBody{},
		}

		res, err := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), blk.Slot)
		require.NoError(t, err)
		require.Equal(t, res.Status, core.HighestProposalVote)
	})
}
