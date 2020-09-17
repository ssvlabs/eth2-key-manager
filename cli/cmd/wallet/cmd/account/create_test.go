package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountCreate(t *testing.T) {
	t.Run("Successfully create account", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--storage=7b226163636f756e7473223a2237623764222c226174744d656d6f7279223a2237623764222c226e6574776f726b223a223664363136393665222c2270726f706f73616c4d656d6f7279223a2237623764222c2277616c6c6574223a2237623232363936343232336132323336363333383333333633353333333732643338333533363632326433343631363436333264363136313337333232643336363533393339333536363631333736363339333933393232326332323639366536343635373834643631373037303635373232323361376237643263323237343739373036353232336132323438343432323764227d",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully create account at specific index", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--storage=7b226163636f756e7473223a2237623764222c226174744d656d6f7279223a2237623764222c226e6574776f726b223a223664363136393665222c2270726f706f73616c4d656d6f7279223a2237623764222c2277616c6c6574223a2237623232363936343232336132323336363333383333333633353333333732643338333533363632326433343631363436333264363136313337333232643336363533393339333536363631333736363339333933393232326332323639366536343635373834643631373037303635373232323361376237643263323237343739373036353232336132323438343432323764227d",
			"--index=5",
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
			"--storage=01213",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to HEX decode seed: encoding/hex: odd length hex string")
	})

	t.Run("Fail to JSON un-marshal", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"create",
			"--seed=b5b0177798165f506de1d46e8e5dd131c708c109800c0e0ce7199aec6572f405",
			"--storage=7b226163636f756e7473223a2237623764222c226174744d656d6f7279223a2237623764222c2270726f706f",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to JSON un-marshal storage: unexpected end of JSON input")
	})
}
