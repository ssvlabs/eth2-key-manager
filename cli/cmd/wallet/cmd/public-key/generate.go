package public_key

import (
	rootcmd "github.com/bloxapp/eth-key-manager/cli/cmd"
	"github.com/bloxapp/eth-key-manager/cli/cmd/wallet/cmd/public-key/flag"
	"github.com/bloxapp/eth-key-manager/cli/cmd/wallet/cmd/public-key/handler"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate public key command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a wallet public key.",
	Long:  `This command generates a public key using seed and index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Generate(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddIndexFlag(generateCmd)
	flag.AddSeedFlag(generateCmd)

	Command.AddCommand(generateCmd)
}
