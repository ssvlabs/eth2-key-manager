package account

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/handler"
)

// voluntaryExitCmd represents the voluntary exit account command.
var voluntaryExitCmd = &cobra.Command{
	Use:   "voluntary-exit",
	Short: "Sign voluntary exit message",
	Long:  `This command signing voluntary exit message using seed or preparing request for signing using key-vault`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.VoluntaryExit(cmd, args)
	},
}

func init() {
	// Define flags for the command.
	rootcmd.AddNetworkFlag(voluntaryExitCmd)
	rootcmd.AddSeedFlag(voluntaryExitCmd)
	rootcmd.AddIndexFlag(voluntaryExitCmd)
	rootcmd.AddResponseTypeFlag(voluntaryExitCmd)
	flag.AddCurrentForkVersionFlag(voluntaryExitCmd)
	flag.AddValidatorPublicKeyFlag(voluntaryExitCmd)
	flag.AddValidatorIndexFlag(voluntaryExitCmd)
	flag.AddEpochFlag(voluntaryExitCmd)

	Command.AddCommand(voluntaryExitCmd)
}
