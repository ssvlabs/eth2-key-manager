package handler

import (
	"encoding/hex"
	"sort"

	"github.com/bloxapp/KeyVault/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"
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

	var accounts []map[string]string
	for account := range wallet.Accounts() {
		accObj := map[string]string{
			"validationPubKey": hex.EncodeToString(account.ValidatorPublicKey().Marshal()),
			"basePath":         account.BasePath(),
		}
		accounts = append(accounts, accObj)
	}

	if len(accounts) == 0 {
		h.printer.Text(storageFlagValue)
		return nil
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i]["basePath"] > accounts[j]["basePath"]
	})

	err = wallet.DeleteAccountByPublicKey(accounts[0]["validationPubKey"])
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
