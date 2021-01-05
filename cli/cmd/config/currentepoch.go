package seed

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/config/handler"
)

// currentEpochCmd represents the command to get the current epoch.
var currentEpochCmd = &cobra.Command{
	Use:   "current-epoch",
	Short: "Prints the current epoch.",
	Long:  `This command prints the current epoch based on the network.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		network, err := rootcmd.GetNetworkFlagValue(cmd)
		if err != nil {
			return err
		}

		handler := handler.New(rootcmd.ResultPrinter, network)
		return handler.CurrentEpoch(cmd, args)
	},
}

func init() {
	Command.AddCommand(currentEpochCmd)
}
