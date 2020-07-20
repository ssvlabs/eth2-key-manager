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
	cliflag.AddPersistentBoolFlag(c, mnemonicFlag, false, "Description here...", false)
}

// GetMnemonicFlagValue gets the mnemonic flag from the command
func GetMnemonicFlagValue(c *cobra.Command) (bool, error) {
	return c.Flags().GetBool(mnemonicFlag)
}
