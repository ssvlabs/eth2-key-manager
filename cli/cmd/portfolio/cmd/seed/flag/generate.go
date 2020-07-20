package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/KeyVault/cli/util/cliflag"
)

// Flag names.
const (
	mnemonicFlag = "mnemonic"
	seedFlag = "seed"
)

// AddMnemonicFlag adds the mnemonic flag to the command
func AddMnemonicFlag(c *cobra.Command) {
	cliflag.AddPersistentBoolFlag(c, mnemonicFlag, false, "Generate mnemonic phrase", false)
}

// AddSeedFlag adds the seed flag to the command
func AddSeedFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, seedFlag, "", "Seed to mnemonic phrase", false)
}

// GetMnemonicFlagValue gets the mnemonic flag from the command
func GetMnemonicFlagValue(c *cobra.Command) (bool, error) {
	return c.Flags().GetBool(mnemonicFlag)
}

// GetSeedFlagValue gets the seed flag from the command
func GetSeedFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(seedFlag)
}
