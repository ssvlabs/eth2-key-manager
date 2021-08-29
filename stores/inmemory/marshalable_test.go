package inmemory

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/wallets"
	"github.com/bloxapp/eth2-key-manager/wallets/hd"
	"github.com/bloxapp/eth2-key-manager/wallets/nd"
)

func TestMarshalingWallet(t *testing.T) {
	err := core.InitBLS()
	require.NoError(t, err)

	store := NewInMemStore(core.MainNetwork)

	// setup wallet
	walletCtx := &core.WalletContext{Storage: store}
	wallet := nd.NewWallet(walletCtx)
	k, err := core.NewHDKeyFromPrivateKey(_byteArray("5470813f7deef638dc531188ca89e36976d536f680e89849cd9077fd096e20bc"), "")
	require.NoError(t, err)
	account := wallets.NewValidatorAccount("", k, k.PublicKey().Serialize(), "", walletCtx)
	require.NoError(t, err)
	wallet.AddValidatorAccount(account)
	err = store.SaveWallet(wallet)
	require.NoError(t, err)

	// marshal
	byts, err := json.Marshal(store)
	require.NoError(t, err)

	// un-marshal
	var store2 InMemStore
	require.NoError(t, json.Unmarshal(byts, &store2))

	// verify
	t.Run("verify wallet", func(t *testing.T) {
		wallet2, err := store2.OpenWallet()
		require.NoError(t, err)
		require.Equal(t, wallet.ID().String(), wallet2.ID().String())
	})
	t.Run("verify acc", func(t *testing.T) {
		wallet2, err := store2.OpenWallet()
		require.NoError(t, err)
		acc2, err := wallet2.AccountByPublicKey("a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972")
		require.NoError(t, err)
		require.Equal(t, account.ID().String(), acc2.ID().String())
	})
}

func TestMarshaling(t *testing.T) {
	err := core.InitBLS()
	require.NoError(t, err)

	store := NewInMemStore(core.MainNetwork)

	// wallet
	wallet := hd.NewWallet(&core.WalletContext{Storage: store})
	err = store.SaveWallet(wallet)
	require.NoError(t, err)

	// account
	acc, err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), nil)
	require.NoError(t, err)
	err = store.SaveAccount(acc)
	require.NoError(t, err)

	// attestation
	att := &eth.AttestationData{
		Slot:            1,
		CommitteeIndex:  1,
		BeaconBlockRoot: []byte("A"),
		Source: &eth.Checkpoint{
			Epoch: 1,
			Root:  []byte("A"),
		},
		Target: &eth.Checkpoint{
			Epoch: 2,
			Root:  []byte("A"),
		},
	}
	store.SaveHighestAttestation(acc.ValidatorPublicKey(), att)

	// proposal
	prop := &eth.BeaconBlock{
		Slot:          1,
		ProposerIndex: 1,
		ParentRoot:    []byte("A"),
		StateRoot:     []byte("A"),
		Body:          &eth.BeaconBlockBody{},
	}
	store.SaveHighestProposal(acc.ValidatorPublicKey(), prop)

	// marshal
	byts, err := json.Marshal(store)
	require.NoError(t, err)

	fmt.Printf("%s\n", hex.EncodeToString(byts))

	// un-marshal
	var store2 InMemStore
	require.NoError(t, json.Unmarshal(byts, &store2))

	// verify
	t.Run("verify wallet", func(t *testing.T) {
		wallet2, err := store2.OpenWallet()
		require.NoError(t, err)
		require.Equal(t, wallet.ID().String(), wallet2.ID().String())
	})
	t.Run("verify acc", func(t *testing.T) {
		wallet2, err := store2.OpenWallet()
		require.NoError(t, err)
		acc2, err := wallet2.AccountByPublicKey("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")
		require.NoError(t, err)
		require.Equal(t, acc.ID().String(), acc2.ID().String())
	})
	t.Run("verify attestation", func(t *testing.T) {
		att2 := store.RetrieveHighestAttestation(acc.ValidatorPublicKey())
		require.Equal(t, att.BeaconBlockRoot, att2.BeaconBlockRoot)
	})
	t.Run("verify proposal", func(t *testing.T) {
		prop2 := store.RetrieveHighestProposal(acc.ValidatorPublicKey())
		require.Equal(t, prop.StateRoot, prop2.StateRoot)
	})
}
