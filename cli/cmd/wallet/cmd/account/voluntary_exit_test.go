package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/util/printer"
)

func TestAccountVoluntaryExit(t *testing.T) {
	t.Run("Successfully sign voluntary exit", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--response-type=object",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--current-fork-version=0x02001020",
			"--index=0",
			"--validator-index=273230",
			"--validator-public-key=b2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c",
			"--epoch=183797",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully prepare sign voluntary exit request for key-vault", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--current-fork-version=0x02001020",
			"--index=0",
			"--validator-index=273230",
			"--validator-public-key=b2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c",
			"--epoch=183797",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Invalid current fork version length", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--current-fork-version=0x020000",
			"--index=1",
			"--validator-index=1",
			"--validator-public-key=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c",
			"--epoch=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect voluntary exit flags: failed to retrieve the current fork version flag value: invalid length for current fork version")
	})

	t.Run("Invalid validator public key", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--current-fork-version=0x02001020",
			"--index=1",
			"--validator-index=1",
			"--validator-public-key=0x2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c",
			"--epoch=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect voluntary exit flags: failed to parse validator public key: invalid validator public key supplied: encoding/hex: odd length hex string")
	})

	t.Run("Seed flag is required for object response type", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--response-type=object",
			"--seed=",
			"--current-fork-version=0x02001020",
			"--index=0",
			"--validator-index=273230",
			"--validator-public-key=b2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c",
			"--epoch=183797",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect voluntary exit flags: seed flag is required for object response type")
	})
}
