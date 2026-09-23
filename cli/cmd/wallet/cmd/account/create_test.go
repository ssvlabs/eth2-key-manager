package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/util/printer"
)

func TestAccountCreate(t *testing.T) {
	t.Run("Successfully create account at specific index and return as object (prater)", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--response-type=object",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("Successfully create account at specific index and return as object (mainnet)", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--response-type=object",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=mainnet",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("no network flag", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--response-type=object",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=not_known",
		})
		err := cmd.RootCmd.Execute()
		require.Empty(t, output.String())
		require.EqualError(t, err, "failed to collect account flags: failed to retrieve the network flag value: unknown network not_known")
	})

	t.Run("Successfully create account at specific index and return as storage", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=0",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("Successfully create accounts till specific index and return as objects", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
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
			"--highest-source=1,2,3,4,5,6",
			"--highest-target=2,3,4,5,6,7",
			"--highest-proposal=2,3,4,5,6,7",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("Successfully create accounts till specific index and return as storage", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--accumulate=true",
			"--highest-source=1,2,3,4,5,6",
			"--highest-target=2,3,4,5,6,7",
			"--highest-proposal=2,3,4,5,6,7",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("Fail to HEX decode seed", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=01213",
			"--index=1",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.EqualError(t, err, "failed to collect account flags: failed to HEX decode seed: encoding/hex: odd length hex string")
	})

	t.Run("highest sources invalid (accumulate false)", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=1",
			"--highest-source=1,2,3,4,5",
			"--highest-target=2,3,4,5,6",
			"--highest-proposal=2,3,4,5,6",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.EqualError(t, err, "failed to collect account flags: highest sources length when the accumulate flag is false needs to be 1: got 5, need 1")
	})

	t.Run("highest proposal invalid (accumulate false)", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=1",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2,3,4,5,6,7",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.EqualError(t, err, "failed to collect account flags: highest proposals length when the accumulate flag is false needs to be 1: got 6, need 1")
	})

	t.Run("highest sources invalid (accumulate true)", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=1",
			"--accumulate=true",
			"--highest-source=1",
			"--highest-target=2,3",
			"--highest-proposal=2,3",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.EqualError(t, err, "failed to collect account flags: highest sources length when the accumulate flag is true needs to be index + 1: got 1, need 2")
	})

	t.Run("Successfully create seedless account at specific index and return as object (prater)", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index=5",
			"--response-type=object",
			"--highest-source=6",
			"--highest-target=6",
			"--highest-proposal=6",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("Successfully create seedless account at specific index and return as storage", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index=0",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("Successfully create 3 seedless accounts from specific index and return as objects", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898,63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989899,63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989890",
			"--index=1",
			"--response-type=object",
			"--highest-proposal=2,3,4",
			"--highest-target=2,3,4",
			"--highest-source=2,3,4",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.NotEmpty(t, output.String())
	})

	t.Run("Missing Highest Values", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.EqualError(t, err, `required flag(s) "highest-proposal", "highest-source", "highest-target" not set`)
	})

	t.Run("highest sources invalid (seedless)", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index=1",
			"--highest-source=1,2",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.EqualError(t, err, "failed to collect account flags: highest sources length for seedless accounts needs to be equal to private keys count: got 2, need 1")
	})

	t.Run("Fail to HEX decode private key", func(t *testing.T) {
		resetFlags(t, cmd.RootCmd)
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--private-key=01213",
			"--index=1",
			"--network=prater",
			"--highest-proposal=2",
			"--highest-target=2",
			"--highest-source=2",
		})
		err := cmd.RootCmd.Execute()
		require.EqualError(t, err, "failed to collect account flags: failed to HEX decode private-key: encoding/hex: odd length hex string")
	})
}
