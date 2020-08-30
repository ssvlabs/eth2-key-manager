package handler

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth-key-manager/cli/cmd/seed/flag"
	"github.com/bloxapp/eth-key-manager/core"
)

// Seed generates a new key-vault seed and prints it.
func (h *Seed) Generate(cmd *cobra.Command, args []string) error {
	var seed []byte

	// Get mnemonic flag.
	mnemonicFlagValue, err := flag.GetMnemonicFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the mnemonic flag value")
	}

	if len(mnemonicFlagValue) > 0 {
		seed, err = core.SeedFromMnemonic(mnemonicFlagValue, "")
		if err != nil {
			return errors.Wrap(err, "failed to retrieve seed from mnemonic")
		}
	} else {
		entropy, err := core.GenerateNewEntropy()
		if err != nil {
			return errors.Wrap(err, "failed to generate entropy")
		}
		seed, err = core.SeedFromEntropy(entropy, "")
		if err != nil {
			return errors.Wrap(err, "failed to generate seed from entropy")
		}
	}

	h.printer.Text(hex.EncodeToString(seed))
	return nil
}
