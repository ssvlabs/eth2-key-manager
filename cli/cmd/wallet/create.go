package wallet

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/handler"
)

// createCmd represents the create wallet command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a wallet.",
	Long:  `This command creates a wallet.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter, rootcmd.Network)
		return handler.Create(cmd, args)
	},
}

func init() {
	Command.AddCommand(createCmd)
}
