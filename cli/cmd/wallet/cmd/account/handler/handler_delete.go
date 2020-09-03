package handler

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
)

// Account deletes a last indexed account and prints the storage.
func (h *Account) Delete(cmd *cobra.Command, args []string) error {
	err := types.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
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

	accounts := wallet.Accounts()
	if len(accounts) == 0 {
		h.printer.Text(storageFlagValue)
		return nil
	}

	err = wallet.DeleteAccountByPublicKey(hex.EncodeToString(accounts[0].ValidatorPublicKey().Marshal()))
	if err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	// marshal storage
	bytes, err := store.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "failed to JSON marshal storage")
	}

	h.printer.Text(hex.EncodeToString(bytes))
	return nil
}
