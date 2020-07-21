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
			"portfolio",
			"seed",
			"generate",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
	})

	t.Run("Successfully Generate Mnemonic", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"portfolio",
			"seed",
			"generate",
			"--mnemonic",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
	})

	t.Run("Successfully Generate Mnemonic from seed", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"portfolio",
			"seed",
			"generate",
			"--mnemonic",
			"--seed=a42b2d973095bb518e45ae5b372dbff9a3aec572ff74b1c8c54749d34b4479eb",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.Equal(t, "picnic floor grape gentle forum potato decorate remind forest ride husband viable depend glance slogan update range economy fade neck cruise peasant tray illegal\n", actualOutput)
	})

	t.Run("Don't generate seed if seed flag is provided", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"portfolio",
			"seed",
			"generate",
			"--mnemonic=false",
			"--seed=a42b2d973095bb518e45ae5b372dbff9a3aec572ff74b1c8c54749d34b4479eb",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.Equal(t, "a42b2d973095bb518e45ae5b372dbff9a3aec572ff74b1c8c54749d34b4479eb\n", actualOutput)
	})
}
