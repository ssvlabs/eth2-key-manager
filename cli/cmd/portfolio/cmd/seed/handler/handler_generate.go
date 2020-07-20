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
	// Get mnemonic flag.
	mnemonic, err := flag.GetMnemonicFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the mnemonic flag value")
	}

	h.printer.JSON(mnemonic)

	seed, err := KeyVault.GenerateNewSeed()
	if err != nil {
		return err
	}
	h.printer.Text(hex.EncodeToString(seed))
	return nil
}
