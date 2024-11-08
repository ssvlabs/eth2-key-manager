package mnemonic_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/util/printer"
)

func TestMnemonicGenerate(t *testing.T) {
	t.Run("Successfully Generate Mnemonic", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"mnemonic",
			"generate",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
	})
}
