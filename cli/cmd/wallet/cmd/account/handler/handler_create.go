package handler

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/KeyVault/cli/cmd/wallet/cmd/account/flag"
)

// Account creates a new wallet account and prints the storage.
func (h *Account) Create(cmd *cobra.Command, args []string) error {
	types.InitBLS()

	// Get seed flag.
	nameFlagValue, err := flag.GetNameFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the name flag value")
	}

	// Get seed flag.
	seedFlagValue, err := flag.GetSeedFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	seedBytes, err := hex.DecodeString(seedFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode seed")
	}

	// Get storage flag.
	storageFlagValue, err := flag.GetStorageFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the storage flag value")
	}

	storageBytes, err := hex.DecodeString(storageFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode storage")
	}

	var store in_memory.InMemStore
	err = store.UnmarshalJSON(storageBytes)
	if err != nil {
		return errors.Wrap(err, "failed to JSON un-marshal storage")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open key wallet")
	}

	_, err = wallet.CreateValidatorAccount(seedBytes, nameFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to create validator account")
	}

	// marshal storage
	bytes, err := store.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "failed to JSON marshal storage")
	}

	h.printer.Text(hex.EncodeToString(bytes))
	return nil
}
