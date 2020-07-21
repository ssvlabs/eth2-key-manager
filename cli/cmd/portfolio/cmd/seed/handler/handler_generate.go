package handler

import (
	"encoding/hex"

	"github.com/bloxapp/KeyVault"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/KeyVault/cli/cmd/portfolio/cmd/seed/flag"
)

// Seed generates a new portfolio seed and prints it.
func (h *Seed) Generate(cmd *cobra.Command, args []string) error {
	var seed []byte
	var mnemonic string

	// Get mnemonic flag.
	mnemonicFlagValue, err := flag.GetMnemonicFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the mnemonic flag value")
	}

	// Get seed flag.
	seedFlagValue, err := flag.GetSeedFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	// Generate seed if seed flag is not provided.
	if len(seedFlagValue) == 0 {
		seed, err = KeyVault.GenerateNewSeed()
		if err != nil {
			return errors.Wrap(err, "failed to generate new seed")
		}
	} else if mnemonicFlagValue {
		seed, err = hex.DecodeString(seedFlagValue)
		if err != nil {
			return errors.Wrap(err, "failed to hex decode the seed flag value")
		}
	}

	// Generate mnemonic
	if mnemonicFlagValue {
		mnemonic, err = KeyVault.SeedToMnemonic(seed)
		if err != nil {
			return errors.Wrap(err, "failed to generate mnemonic from seed")
		}
		h.printer.Text(mnemonic)
	} else if len(seedFlagValue) > 0 {
		h.printer.Text(seedFlagValue)
	} else {
		h.printer.Text(hex.EncodeToString(seed))
	}

	return nil
}
