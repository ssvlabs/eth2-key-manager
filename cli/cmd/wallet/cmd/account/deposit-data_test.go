package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountDepositData(t *testing.T) {
	t.Run("Successfully retrieve deposit-data", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"deposit-data",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=0",
			"--publickey=95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf",
			"--network=pyrmont",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully retrieve deposit-data for launchtest network", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"deposit-data",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=0",
			"--publickey=95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf",
			"--network=pyrmont",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Fail retrieve deposit-data for unmatched index and publickey", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"deposit-data",
			"--seed=0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff",
			"--index=5",
			"--publickey=81fd26fe6e7cdbe1d0d45020050ba94c625f5236bf162b9ad3fca137d9120a0572c6f59b8cc70fae6cd6bb471b673e97",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to get account by public key: account not found")
	})
}
