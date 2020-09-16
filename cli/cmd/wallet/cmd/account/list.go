package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// listCmd represents the list wallet accounts command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallet accounts.",
	Long:  `This command list wallet accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.List(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddStorageFlag(listCmd)

	Command.AddCommand(listCmd)
}
