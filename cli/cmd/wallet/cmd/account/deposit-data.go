package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// depositDataCmd represents the deposit-data account command.
var depositDataCmd = &cobra.Command{
	Use:   "deposit-data",
	Short: "Returns an account deposit-data.",
	Long:  `This command returns an account deposit-data using public key and storage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter, rootcmd.Network)
		return handler.DepositData(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	flag.AddPublicKeyFlag(depositDataCmd)
	flag.AddStorageFlag(depositDataCmd)

	Command.AddCommand(depositDataCmd)
}
