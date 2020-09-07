package validator

import (
	"github.com/spf13/cobra"

	keyvaultcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
)

// Command represents the validator related command.
var Command = &cobra.Command{
	Use:   "validator",
	Short: "Manage validators",
}

func init() {
	keyvaultcmd.RootCmd.AddCommand(Command)
}
