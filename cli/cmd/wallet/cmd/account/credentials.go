package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// credentialsCmd represents the credentials account command.
var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Execute BLS to execution",
	Long:  `This command executing BLS to execution change using seed`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Credentials(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	rootcmd.AddNetworkFlag(credentialsCmd)
	rootcmd.AddSeedFlag(credentialsCmd)
	rootcmd.AddIndexFlag(credentialsCmd)
	rootcmd.AddAccumulateFlag(credentialsCmd)
	flag.AddValidatorIndexFlag(credentialsCmd)
	flag.AddValidatorPublicKeyFlag(credentialsCmd)
	flag.AddWithdrawalCredentialsFlag(credentialsCmd)
	flag.AddToExecutionAddressFlag(credentialsCmd)

	Command.AddCommand(credentialsCmd)
}
