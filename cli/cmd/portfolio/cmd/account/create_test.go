package seed_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/KeyVault/cli/cmd"
	"github.com/bloxapp/KeyVault/cli/util/printer"
)

func TestAccountCreate(t *testing.T) {
	t.Run("Successfully create account", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"portfolio",
			"account",
			"create",
			"--index=0",
			"--seed=a42b2d973095bb518e45ae5b372dbff9a3aec572ff74b1c8c54749d34b4479eb",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		actualOutput := output.String()
		require.Equal(t, "235f794fb6d5c647cee25d23a89d94cc65c6c8fee0bf77685d738a81d703339e\n", actualOutput)
	})

	t.Run("Fail to HEX decode seed", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"portfolio",
			"account",
			"create",
			"--index=0",
			"--seed=01213",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to HEX decode seed: encoding/hex: odd length hex string")
	})
}
