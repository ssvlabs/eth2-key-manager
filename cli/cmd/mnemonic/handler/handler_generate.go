package handler

import (
	"github.com/bloxapp/eth-key-manager/core"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Mnemonic generates a new key-vault mnemonic and prints it.
func (h *Mnemonic) Generate(cmd *cobra.Command, args []string) error {
	// Generate new entropy
	entropy, err := core.GenerateNewEntropy()
	if err != nil {
		return errors.Wrap(err, "failed to generate entropy")
	}

	// Generate mnemonic from entropy
	mnemonic, err := core.EntropyToMnemonic(entropy)
	if err != nil {
		return errors.Wrap(err, "failed to generate mnemonic from entropy")
	}

	h.printer.Text(mnemonic)
	return nil
}
