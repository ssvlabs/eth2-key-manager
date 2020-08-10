package handler

import (
	"encoding/hex"
	"fmt"

	"github.com/bloxapp/KeyVault/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/KeyVault/eth1_deposit"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"
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

	account, err := wallet.AccountByPublicKey(publicKeyFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to get account by public key")
	}

	var accountIndex int
	accountIndex, err = fmt.Sscanf(account.BasePath(), "/%d", &accountIndex)
	if err != nil {
		return errors.Wrap(err, "failed to parse basePath to index")
	}

	withdrawalKey, err := wallet_hd.CreatePrivateKey(seedBytes, wallet_hd.WithdrawalKeyPath, accountIndex)
	if err != nil {
		return errors.Wrap(err, "failed to create withdrawal key")
	}

	depositData, root, err := eth1_deposit.DepositData(account, withdrawalKey, eth1_deposit.MaxEffectiveBalanceInGwei)
	if err != nil {
		return errors.Wrap(err, "failed to get deposit data")
	}

	h.printer.JSON(map[string]interface{}{
		"amount":                depositData.GetAmount(),
		"publicKey":             hex.EncodeToString(depositData.GetPublicKey()),
		"signature":             hex.EncodeToString(depositData.GetSignature()),
		"withdrawalCredentials": hex.EncodeToString(depositData.GetWithdrawalCredentials()),
		"depositDataRoot":       hex.EncodeToString(root[:]),
	})
	return nil
}
