package handler

import (
	"github.com/ssvlabs/eth2-key-manager/core"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Generate generates a new key-vault mnemonic and prints it.
func (h *Mnemonic) Generate(_ *cobra.Command, _ []string) error {
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
