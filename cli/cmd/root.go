package cmd

import (
	"fmt"
	"os"

	"github.com/bloxapp/eth2-key-manager/core"

	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

// ResultPrinter is the printer used to print command results and errors.
var ResultPrinter printer.Printer = printer.NewStandardOutputPrinter()

// Network is the network that could be redefined above.
var Network core.Network = core.TestNetwork

// RootCmd represents the base command when called without any sub-commands.
var RootCmd = &cobra.Command{
	Use:   "keyvault-cli",
	Short: "KeyVault CLI",
	Long:  `keyvault-cli is a CLI for running KeyVault operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		Network, err = GetNetworkFlagValue(cmd)
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	// Register global flags.
	AddNetworkFlag(RootCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
