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
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-index=273230,273407",
			"--epoch=157965",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})
}
