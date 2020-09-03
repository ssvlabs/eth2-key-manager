package mnemonic

import (
	"github.com/spf13/cobra"

	keyvaultcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
)

// Command represents the key-vault mnemonic related command.
var Command = &cobra.Command{
	Use:   "mnemonic",
	Short: "Manage key-vault mnemonic",
}

func init() {
	keyvaultcmd.RootCmd.AddCommand(Command)
}
