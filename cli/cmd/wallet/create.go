package wallet

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/wallet/handler"
)

// createCmd represents the create wallet command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a wallet.",
	Long:  `This command creates a wallet.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		network, err := rootcmd.GetNetworkFlagValue(cmd)
		if err != nil {
			return err
		}

		handler := handler.New(rootcmd.ResultPrinter, network)
		return handler.Create(cmd, args)
	},
}

func init() {
	rootcmd.AddNetworkFlag(createCmd)

	Command.AddCommand(createCmd)
}
