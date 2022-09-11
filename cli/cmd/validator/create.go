package validator

import (
	"io"
	"os"
	"path/filepath"

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
		network, err := rootcmd.GetNetworkFlagValue(cmd)
		if err != nil {
			return err
		}

		handler := handler.New(rootcmd.ResultPrinter, ResultFactory, network)
		return handler.Create(cmd, args)
	},
}

func fileResultFactory(name string) (io.Writer, func(), error) {
	outFile, err := os.Create(filepath.Clean(name + ".zip"))
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
	rootcmd.AddNetworkFlag(createCmd)

	Command.AddCommand(createCmd)
}
