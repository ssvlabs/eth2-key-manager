package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/KeyVault/cli/util/cliflag"
)

// Flag names.
const (
	mnemonicFlag = "mnemonic"
)

// AddMnemonicFlag adds the mnemonic flag to the command
func AddMnemonicFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, mnemonicFlag, "", "Generate seed from mnemonic", false)
}

// GetMnemonicFlagValue gets the mnemonic flag from the command
func GetMnemonicFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(mnemonicFlag)
}
