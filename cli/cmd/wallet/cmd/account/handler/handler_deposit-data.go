package handler

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
)

// Account DepositData generates account deposit-data and prints it.
func (h *Account) DepositData(cmd *cobra.Command, args []string) error {
	err := types.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	// Get public key flag.
	publicKeyFlagValue, err := flag.GetPublicKeyFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the public key flag value")
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

	fmt.Println(hex.EncodeToString(wallet.Accounts()[0].ValidatorPublicKey().Marshal()))

	account, err := wallet.AccountByPublicKey(publicKeyFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to get account by public key")
	}

	depositData, err := account.GetDepositData()
	if err != nil {
		return errors.Wrap(err, "failed to get deposit data")
	}

	err = h.printer.JSON(depositData)
	if err != nil {
		return errors.Wrap(err, "failed to print deposit-data JSON")
	}
	return nil
}
