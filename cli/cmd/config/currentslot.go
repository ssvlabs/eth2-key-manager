package seed

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/config/handler"
)

// currentSlotCmd represents the command to get the current slot.
var currentSlotCmd = &cobra.Command{
	Use:   "current-slot",
	Short: "Prints the current slot.",
	Long:  `This command prints the current slot based on the network.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		network, err := rootcmd.GetNetworkFlagValue(cmd)
		if err != nil {
			return err
		}

		handler := handler.New(rootcmd.ResultPrinter, network)
		return handler.CurrentSlot(cmd, args)
	},
}

func init() {
	Command.AddCommand(currentSlotCmd)
}
