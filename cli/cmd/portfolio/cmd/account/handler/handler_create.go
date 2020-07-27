package handler

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/KeyVault/cli/cmd/portfolio/cmd/account/flag"
)

// Account creates a new portfolio account and prints it's private key.
func (h *Account) Create(cmd *cobra.Command, args []string) error {
	// Get index flag.
	indexFlagValue, err := flag.GetIndexFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the mnemonic flag value")
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

	key, err := KeyVault.CreateAccount(seedBytes, indexFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to create account")
	}

	h.printer.Text(hex.EncodeToString(key))
	return nil
}
