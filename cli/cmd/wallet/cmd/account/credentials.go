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
	flag.AddIndexFlag(credentialsCmd)
	flag.AddSeedFlag(credentialsCmd)
	flag.AddValidatorIndexFlag(credentialsCmd)
	flag.AddValidatorPublicKeyFlag(credentialsCmd)
	flag.AddWithdrawalCredentialsFlag(credentialsCmd)
	flag.AddToExecutionAddressFlag(credentialsCmd)
	rootcmd.AddAccumulateFlag(credentialsCmd)
	rootcmd.AddNetworkFlag(credentialsCmd)

	Command.AddCommand(credentialsCmd)
}
