package publickey

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/wallet/cmd/publickey/handler"
)

// generateCmd represents the generate public key command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a wallet public key.",
	Long:  `This command generates a public key using seed and index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		network, err := rootcmd.GetNetworkFlagValue(cmd)
		if err != nil {
			return err
		}

		handler := handler.New(rootcmd.ResultPrinter, network)
		return handler.Generate(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	rootcmd.AddNetworkFlag(generateCmd)
	rootcmd.AddSeedFlag(generateCmd)
	rootcmd.AddIndexFlag(generateCmd)

	Command.AddCommand(generateCmd)
}
