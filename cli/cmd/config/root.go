package seed

import (
	"github.com/spf13/cobra"

	keyvaultcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd"
)

// Command represents the key-vault seed related command.
var Command = &cobra.Command{
	Use:   "config",
	Short: "Manage key-vault config",
}

func init() {
	keyvaultcmd.AddNetworkFlag(Command)

	keyvaultcmd.RootCmd.AddCommand(Command)
}
