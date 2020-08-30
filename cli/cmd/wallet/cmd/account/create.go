package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth-key-manager/cli/cmd"
	"github.com/bloxapp/eth-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// createCmd represents the create account command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a wallet account.",
	Long:  `This command creates an account using seed and storage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Create(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddNameFlag(createCmd)
	flag.AddSeedFlag(createCmd)
	flag.AddStorageFlag(createCmd)

	Command.AddCommand(createCmd)
}
