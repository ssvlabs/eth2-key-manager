package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountVoluntaryExit(t *testing.T) {
	t.Run("Successfully handle voluntary exit", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--seed=94227e2b758770d5703c7056d060d2bcdab49364bf14c917d43c71c4e619164c0e53e577cebfa8587e6febd36c3dd02eb44f8a0916b2aaabdf902a6bfe99caeb",
			"--index=0",
			//"--accumulate=true",
			//"--validator-index=273230,273407",
			"--validator-index=0",
			"--epoch=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})
}
