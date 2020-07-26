package seed

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/KeyVault/cli/cmd"
	"github.com/bloxapp/KeyVault/cli/cmd/portfolio/cmd/account/flag"
	"github.com/bloxapp/KeyVault/cli/cmd/portfolio/cmd/account/handler"
)

// generateCmd represents the create account command.
var generateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a wallet account.",
	Long:  `This command creates an account using seed and index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Create(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddIndexFlag(generateCmd)
	flag.AddSeedFlag(generateCmd)

	Command.AddCommand(generateCmd)
}
