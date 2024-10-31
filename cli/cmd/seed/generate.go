package seed

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/seed/flag"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/seed/handler"
)

// generateCmd represents the generate seed command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a key-vault seed.",
	Long:  `This command generates a random seed of the key-vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Generate(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddMnemonicFlag(generateCmd)

	Command.AddCommand(generateCmd)
}
