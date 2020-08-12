package handler

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"
	"sort"
)

// Account list wallet accounts and prints the accounts.
func (h *Account) List(cmd *cobra.Command, args []string) error {
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
	for a := range wallet.Accounts() {
		accObj := map[string]string{
			"id":               a.ID().String(),
			"name":             a.Name(),
			"validationPubKey": hex.EncodeToString(a.ValidatorPublicKey().Marshal()),
			"withdrawalPubKey": hex.EncodeToString(a.WithdrawalPublicKey().Marshal()),
		}
		accounts = append(accounts, accObj)
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i]["name"] > accounts[j]["name"]
	})

	err = h.printer.JSON(accounts)
	if err != nil {
		return errors.Wrap(err, "failed to print accounts JSON")
	}
	return nil
}
