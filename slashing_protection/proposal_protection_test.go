package slashing_protection

import (
	"encoding/hex"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/stretchr/testify/require"

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

func setupProposal(updateHighestProposal bool) (core.SlashingProtector, []core.ValidatorAccount, error) {
	if err := core.InitBLS(); err != nil { // very important!
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

	if updateHighestProposal {
		protector.UpdateHighestProposal(account1.ValidatorPublicKey(), &eth.BeaconBlock{
			Slot:          100,
			ProposerIndex: 2,
			ParentRoot:    []byte("A"),
			StateRoot:     []byte("A"),
			Body:          &eth.BeaconBlockBody{},
		})
	}

	return protector, []core.ValidatorAccount{account1, account2}, nil
}

func TestProposalProtection(t *testing.T) {
	t.Run("New proposal, should not slash", func(t *testing.T) {
		protector, accounts, err := setupProposal(true)
		require.NoError(t, err)
		res, err := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), &eth.BeaconBlock{
			Slot:          101,
			ProposerIndex: 2,
			ParentRoot:    []byte("Z"),
			StateRoot:     []byte("Z"),
			Body:          &eth.BeaconBlockBody{},
		})
		require.NoError(t, err)
		require.Equal(t, res.Status, core.ValidProposal)
	})

	t.Run("No highest proposal db, should error", func(t *testing.T) {
		protector, accounts, err := setupProposal(false)
		require.NoError(t, err)
		res, err := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), &eth.BeaconBlock{
			Slot:          99,
			ProposerIndex: 2,
			ParentRoot:    []byte("Z"),
			StateRoot:     []byte("Z"),
			Body:          &eth.BeaconBlockBody{},
		})
		require.EqualError(t, err, "highest proposal data is nil, can't determine if proposal is slashable")
		require.Nil(t, res)
	})

	t.Run("Lower than highest proposal db, should error", func(t *testing.T) {
		protector, accounts, err := setupProposal(true)
		require.NoError(t, err)
		res, err := protector.IsSlashableProposal(accounts[0].ValidatorPublicKey(), &eth.BeaconBlock{
			Slot:          99,
			ProposerIndex: 2,
			ParentRoot:    []byte("Z"),
			StateRoot:     []byte("Z"),
			Body:          &eth.BeaconBlockBody{},
		})
		require.NoError(t, err)
		require.Equal(t, res.Status, core.HighestProposalVote)
	})
}
