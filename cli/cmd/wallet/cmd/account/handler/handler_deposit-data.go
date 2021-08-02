package handler

import (
	"encoding/hex"

	"github.com/bloxapp/eth2-key-manager/core"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
)

// DepositData generates account deposit-data and prints it.
func (h *Account) DepositData(cmd *cobra.Command, _ []string, seedless bool) error {
	err := core.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	var seedOrPrivateKeyBytes []byte
	var seedOrPrivateKeyFlagValue string

	if seedless == true {
		// Get private key flag value
		seedOrPrivateKeyFlagValue, err = flag.GetPrivateKeyFlagValue(cmd)
	} else {
		// Get seed flag value
		seedOrPrivateKeyFlagValue, err = flag.GetSeedFlagValue(cmd)
	}
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	seedOrPrivateKeyBytes, err = hex.DecodeString(seedOrPrivateKeyFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode seed")
	}

	// Get network flag
	network, err := rootcmd.GetNetworkFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the network flag value")
	}

	// Get public key flag.
	publicKeyFlagValue, err := flag.GetPublicKeyFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the public key flag value")
	}

	// TODO get rid of network
	store := inmemory.NewInMemStore(network)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	_, err = eth2keymanager.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	// Get index flag.
	indexFlagValue, err := flag.GetIndexFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the index flag value")
	}

	if seedless == true {
		_, err = wallet.CreateValidatorAccountFromPrivateKey(seedOrPrivateKeyBytes, &indexFlagValue)
	} else {
		_, err = wallet.CreateValidatorAccount(seedOrPrivateKeyBytes, &indexFlagValue)
	}
	if err != nil {
		return errors.Wrap(err, "failed to create validator account")
	}

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
