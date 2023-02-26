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
	Short: "Execute voluntary exit",
	Long:  `This command executing voluntary exit using seed`,
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
	rootcmd.AddAccumulateFlag(voluntaryExitCmd)
	flag.AddCurrentForkVersionFlag(voluntaryExitCmd)
	flag.AddValidatorIndexFlag(voluntaryExitCmd)
	flag.AddEpochFlag(voluntaryExitCmd)

	Command.AddCommand(voluntaryExitCmd)
}
