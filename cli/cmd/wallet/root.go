package wallet

import (
	"github.com/spf13/cobra"

	keyvaultcmd "github.com/bloxapp/KeyVault/cli/cmd"
)

// Command represents the wallet related command.
var Command = &cobra.Command{
	Use:   "wallet",
	Short: "Manage wallet",
}

func init() {
	keyvaultcmd.RootCmd.AddCommand(Command)
}
