package handler

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault"
	"github.com/spf13/cobra"
)

// Seed generates a new portfolio seed and prints it.
func (h *Seed) Generate(cmd *cobra.Command, args []string) error {
	seed, err := KeyVault.GenerateNewSeed()
	if err != nil {
		return err
	}
	h.printer.Text(hex.EncodeToString(seed))
	return nil
}
