package handler

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth-key-manager/stores/in_memory"
)

// Account creates a new wallet account and prints the storage.
func (h *Account) Create(cmd *cobra.Command, args []string) error {
	err := types.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	var indexPointer *int
	// check if indexFlag was assigned
	if cmd.Flags().Changed(flag.GetIndexFlagName()) {
		// Get index flag.
		indexFlagValue, err := flag.GetIndexFlagValue(cmd)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve the index flag value")
		}
		indexPointer = &indexFlagValue
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
		return errors.Wrap(err, "failed to open wallet")
	}

	_, err = wallet.CreateValidatorAccount(seedBytes, indexPointer)
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
