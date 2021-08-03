package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountCreateSeedless(t *testing.T) {
	t.Run("Successfully create seedless account at specific index-from and return as object (prater)", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index-from=5",
			"--response-type=object",
			"--highest-source=1,2,3,4,5,6",
			"--highest-target=1,2,3,4,5,6",
			"--highest-proposal=1,2,3,4,5,6",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully create seedless account at specific index-from and return as object (mainnet)", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index-from=11",
			"--response-type=object",
			"--highest-source=1,2,3,4,5,6,7,8,9,10,11,12",
			"--highest-target=1,2,3,4,5,6,7,8,9,10,11,12",
			"--highest-proposal=1,2,3,4,5,6,7,8,9,10,11,12",
			"--network=mainnet",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("no network flag for seedless account", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index-from=0",
			"--response-type=object",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=not_known",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.EqualError(t, err, "failed to network: unknown network")
	})

	t.Run("Successfully create seedless account at specific index-from and return as storage", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index-from=0",
			"--highest-source=1",
			"--highest-target=2",
			"--highest-proposal=2",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully create 3 seedless accounts from specific index-from and return as objects", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898,63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989899,63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989890",
			"--index-from=1",
			"--response-type=object",
			"--highest-proposal=1,2,3,4",
			"--highest-target=1,2,3,4",
			"--highest-source=1,2,3,4",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Missing Highest Values", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f08989898",
			"--index-from=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "highest sources length for seedless accounts need to be equal to <index-from> + <private keys count>")
	})

	t.Run("Fail to HEX decode private key", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create-seedless",
			"--private-key=01213",
			"--index-from=1",
			"--network=prater",
			"--highest-proposal=1,2",
			"--highest-target=1,2",
			"--highest-source=1,2",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to convert private key string to bytes: encoding/hex: odd length hex string")
	})
}
