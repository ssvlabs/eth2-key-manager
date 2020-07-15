package seed

import (
	"github.com/spf13/cobra"

	portfoliocmd "github.com/bloxapp/KeyVault/cli/cmd/portfolio"
)

// Command represents the portfolio seed related command.
var Command = &cobra.Command{
	Use:   "seed",
	Short: "Manage portfolio seed",
}

func init() {
	portfoliocmd.Command.AddCommand(Command)
}
