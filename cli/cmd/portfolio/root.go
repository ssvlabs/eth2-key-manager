package portfolio

import (
	"github.com/spf13/cobra"

	keyvaultcmd "github.com/bloxapp/KeyVault/cli/cmd"
)

// Command represents the portfolio-related command.
var Command = &cobra.Command{
	Use:   "portfolio",
	Short: "Manage portfolio",
}

func init() {
	keyvaultcmd.RootCmd.AddCommand(Command)
}
