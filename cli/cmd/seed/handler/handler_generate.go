package handler

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault/core"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/KeyVault/cli/cmd/seed/flag"
)

// Seed generates a new key-vault seed and prints it.
func (h *Seed) Generate(cmd *cobra.Command, args []string) error {
	var seed []byte
	var entropy []byte
	var mnemonic string

	// Get mnemonic flag.
	mnemonicFlagValue, err := flag.GetMnemonicFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the mnemonic flag value")
	}

	// Get seed flag.
	entropyFlagValue, err := flag.GetEntropyFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	// Generate seed if seed flag is not provided.
	if len(entropyFlagValue) == 0 {
		entropy, err = core.GenerateNewEntropy()
		if err != nil {
			return errors.Wrap(err, "failed to generate entropy")
		}

		seed, err = core.SeedFromEntropy(entropy, "")
		if err != nil {
			return errors.Wrap(err, "failed to generate new seed")
		}
	} else if mnemonicFlagValue {
		entropy, err = hex.DecodeString(entropyFlagValue)
		if err != nil {
			return errors.Wrap(err, "failed to hex decode the seed flag value")
		}
	}

	// Generate mnemonic
	if mnemonicFlagValue {
		mnemonic, err = core.EntropyToMnemonic(entropy)
		if err != nil {
			return errors.Wrap(err, "failed to generate mnemonic from seed")
		}
		h.printer.Text(mnemonic)
	} else if len(entropyFlagValue) > 0 {
		h.printer.Text(entropyFlagValue)
	} else {
		h.printer.Text(hex.EncodeToString(seed))
	}

	return nil
}
