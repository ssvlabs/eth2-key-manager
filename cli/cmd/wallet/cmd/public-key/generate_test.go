package public_key_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth-key-manager/cli/cmd"
	"github.com/bloxapp/eth-key-manager/cli/util/printer"
)

func TestPublicKeyGenerate(t *testing.T) {
	t.Run("Successfully generate public-key", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"public-key",
			"generate",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=4",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Fail to generate public-key with negative index", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"public-key",
			"generate",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=-1",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "provided index is negative")
	})
}
