package seed

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/KeyVault/cli/cmd"
	"github.com/bloxapp/KeyVault/cli/cmd/portfolio/cmd/seed/handler"
)

// generateCmd represents the generate seed command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a portfolio seed.",
	Long:  `This command generates a random seed of the current portfolio.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Generate(cmd, args)
	},
}

func init() {
	Command.AddCommand(generateCmd)
}
