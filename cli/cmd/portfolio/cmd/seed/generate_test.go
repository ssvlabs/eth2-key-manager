package seed_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/KeyVault/cli/cmd"
	"github.com/bloxapp/KeyVault/cli/util/printer"
)

func TestSeedGenerate(t *testing.T) {
	var output bytes.Buffer
	cmd.ResultPrinter = printer.New(&output)

	// Execute
	cmd.RootCmd.SetArgs([]string{
		"portfolio",
		"seed",
		"generate",
	})
	err := cmd.RootCmd.Execute()

	// Result
	require.NoError(t, err)

	actualOutput := output.String()
	require.Equal(t, "TODO: Implement seed generating\n", actualOutput)
}
