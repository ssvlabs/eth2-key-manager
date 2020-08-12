package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/KeyVault/cli/cmd"
	"github.com/bloxapp/KeyVault/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/KeyVault/cli/cmd/wallet/cmd/account/handler"
)

// deleteCmd represents the delete account command.
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes last indexed account.",
	Long:  `This command deletes last indexed account in the wallet using the storage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Delete(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddStorageFlag(deleteCmd)

	Command.AddCommand(deleteCmd)
}
