package public_key

import (
	"github.com/spf13/cobra"

	walletcmd "github.com/bloxapp/eth2-key-manager/cli/cmd/wallet"
)

// Command represents the wallet public-key related command.
var Command = &cobra.Command{
	Use:   "public-key",
	Short: "Manage wallet public-keys",
}

func init() {
	walletcmd.Command.AddCommand(Command)
}
