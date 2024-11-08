package publickey

import (
	"github.com/spf13/cobra"

	walletcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd/wallet"
)

// Command represents the wallet publickey related command.
var Command = &cobra.Command{
	Use:   "publickey",
	Short: "Manage wallet public-keys",
}

func init() {
	walletcmd.Command.AddCommand(Command)
}
