package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/KeyVault/cli/util/cliflag"
)

// Flag names.
const (
	mnemonicFlag = "mnemonic"
	entropyFlag  = "entropy"
)

// AddMnemonicFlag adds the mnemonic flag to the command
func AddMnemonicFlag(c *cobra.Command) {
	cliflag.AddPersistentBoolFlag(c, mnemonicFlag, false, "Generate mnemonic phrase", false)
}

// GetMnemonicFlagValue gets the mnemonic flag from the command
func GetMnemonicFlagValue(c *cobra.Command) (bool, error) {
	return c.Flags().GetBool(mnemonicFlag)
}

// AddEntropyFlag adds the seed flag to the command
func AddEntropyFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, entropyFlag, "", "Seed to mnemonic phrase", false)
}

// GetEntropyFlagValue gets the seed flag from the command
func GetEntropyFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(entropyFlag)
}
