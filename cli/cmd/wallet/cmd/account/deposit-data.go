package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// depositDataCmd represents the deposit-data account command.
var depositDataCmd = &cobra.Command{
	Use:   "deposit-data",
	Short: "Returns an account deposit-data.",
	Long:  `This command returns an account deposit-data using public key and storage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.DepositData(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	rootcmd.AddNetworkFlag(depositDataCmd)
	rootcmd.AddSeedFlag(depositDataCmd)
	rootcmd.AddIndexFlag(depositDataCmd)
	flag.AddPublicKeyFlag(depositDataCmd)

	Command.AddCommand(depositDataCmd)
}
