package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountVoluntaryExit(t *testing.T) {
	t.Run("Successfully handle voluntary exit at specific index", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--current-fork-version=0x02000000",
			"--index=0",
			"--validator-index=0",
			"--epoch=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully handle accumulated voluntary exit", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--current-fork-version=0x02000000",
			"--index=1",
			"--accumulate=true",
			"--validator-index=273230,273407",
			"--epoch=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Invalid length for current fork version", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--current-fork-version=0x020000",
			"--index=1",
			"--accumulate=true",
			"--validator-index=273230,273407",
			"--epoch=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect voluntary exit flags: failed to retrieve the current fork version flag value: invalid length for current fork version")
	})

	t.Run("Only one validator can be specified if accumulate is false", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"voluntary-exit",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--current-fork-version=0x02000000",
			"--index=0",
			"--accumulate=false",
			"--validator-index=0,2",
			"--epoch=1",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect voluntary exit flags: only one validator can be specified if accumulate is false")
	})
}
