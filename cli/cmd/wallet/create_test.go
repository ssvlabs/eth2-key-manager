package wallet_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountCreate(t *testing.T) {
	t.Run("Successfully create wallet (prater)", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"create",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
	})

	t.Run("Successfully create wallet (mainnet)", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"create",
			"--network=mainnet",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
	})
}
