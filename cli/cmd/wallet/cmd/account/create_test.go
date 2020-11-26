package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountCreate(t *testing.T) {
	t.Run("Successfully create account at specific index and return as object", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--response-type=object",
			"--minimal-slashing-data=" + "7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a6e756c6c2c22736f75726365223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d7d",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully create account at specific index and return as storage", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=0",
			"--minimal-slashing-data=" + "7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a6e756c6c2c22736f75726365223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d7d",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully create accounts till specific index and return as objects", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--accumulate=true",
			"--response-type=object",
			"--minimal-slashing-data=" + "7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a6e756c6c2c22736f75726365223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d7d",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully create accounts till specific index and return as storage", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--accumulate=true",
			"--minimal-slashing-data=" + "7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a6e756c6c2c22736f75726365223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a6e756c6c7d7d",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Fail to HEX decode seed", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=01213",
			"--index=1",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to HEX decode seed: encoding/hex: odd length hex string")
	})
}
