package seed

import (
	"github.com/spf13/cobra"

	portfoliocmd "github.com/bloxapp/KeyVault/cli/cmd/portfolio"
)

// Command represents the portfolio account related command.
var Command = &cobra.Command{
	Use:   "account",
	Short: "Manage portfolio account",
}

func init() {
	portfoliocmd.Command.AddCommand(Command)
}
