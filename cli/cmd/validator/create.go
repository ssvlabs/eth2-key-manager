package validator

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/validator/flag"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/validator/handler"
)

// ResultFactory is the validation creation result factory
var ResultFactory = fileResultFactory

// createCmd represents the create validator command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates validator(s).",
	Long:  `This command creates validator(s).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := handler.New(rootcmd.ResultPrinter, ResultFactory)
		return handler.Create(cmd, args)
	},
}

func fileResultFactory(name string) (io.Writer, func(), error) {
	outFile, err := os.Create(name + ".zip")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create zip file")
	}
	return outFile, func() {
		outFile.Close()
	}, nil
}

func init() {
	// Define flags for the command.
	flag.AddSeedsCountFlag(createCmd)
	flag.AddValidatorsPerSeedFlag(createCmd)
	flag.AddWalletAddressFlag(createCmd)
	flag.AddWalletPrivateKeyFlag(createCmd)
	flag.AddWeb3AddrFlag(createCmd)

	Command.AddCommand(createCmd)
}
