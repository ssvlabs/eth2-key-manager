package seed

import (
	"github.com/spf13/cobra"

	keyvaultcmd "github.com/bloxapp/eth-key-manager/cli/cmd"
)

// Command represents the key-vault seed related command.
var Command = &cobra.Command{
	Use:   "seed",
	Short: "Manage key-vault seed",
}

func init() {
	keyvaultcmd.RootCmd.AddCommand(Command)
}
