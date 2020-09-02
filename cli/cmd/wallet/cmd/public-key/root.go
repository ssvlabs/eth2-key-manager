package public_key

import (
	walletcmd "github.com/bloxapp/eth-key-manager/cli/cmd/wallet"
	"github.com/spf13/cobra"
)

// Command represents the wallet public-key related command.
var Command = &cobra.Command{
	Use:   "public-key",
	Short: "Manage wallet public-keys",
}

func init() {
	walletcmd.Command.AddCommand(Command)
}
