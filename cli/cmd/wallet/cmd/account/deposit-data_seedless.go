package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// depositDataCmd represents the deposit-data account command.
var depositDataSeedlessCmd = &cobra.Command{
	Use:   "deposit-data-seedless",
	Short: "Returns an account deposit data.",
	Long:  `This command returns an account deposit data using public key and private key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.DepositData(cmd, args, true)
	},
}

func init() {
	// Define flags for the command.
	flag.AddIndexFlag(depositDataSeedlessCmd)
	flag.AddPrivateKeyFlag(depositDataSeedlessCmd)
	flag.AddPublicKeyFlag(depositDataSeedlessCmd)
	rootcmd.AddNetworkFlag(depositDataSeedlessCmd)

	Command.AddCommand(depositDataSeedlessCmd)
}
