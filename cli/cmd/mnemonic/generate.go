package mnemonic

import (
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth-key-manager/cli/cmd"
	"github.com/bloxapp/eth-key-manager/cli/cmd/mnemonic/handler"
)

// generateCmd represents the generate mnemonic command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a key-vault mnemonic.",
	Long:  `This command generates a random mnemonic of the key-vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter)
		return handler.Generate(cmd, args)
	},
}

func init() {
	Command.AddCommand(generateCmd)
}
