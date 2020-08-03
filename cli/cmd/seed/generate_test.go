package seed_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/KeyVault/cli/cmd"
	"github.com/bloxapp/KeyVault/cli/util/printer"
)

func TestSeedGenerate(t *testing.T) {
	t.Run("Successfully Generate Seed", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"seed",
			"generate",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
	})

	t.Run("Successfully Generate Seed from Mnemonic", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"seed",
			"generate",
			"--mnemonic=bone dawn produce network shock transfer magic moment dignity grunt must doll combine olive expose artwork wool wrestle pitch range leg install flip coffee",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.Equal(t, "847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3\n", actualOutput)
	})
}
