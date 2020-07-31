package account

import (
	"github.com/spf13/cobra"

	walletcmd "github.com/bloxapp/KeyVault/cli/cmd/wallet"
)

// Command represents the portfolio account related command.
var Command = &cobra.Command{
	Use:   "account",
	Short: "Manage wallet account",
}

func init() {
	walletcmd.Command.AddCommand(Command)
}
