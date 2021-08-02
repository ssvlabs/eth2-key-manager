package account_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestAccountDepositDataSeedless(t *testing.T) {
	t.Run("Successfully retrieve deposit data in seedless mode", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"deposit-data-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f00633d29",
			"--index=0",
			"--publickey=a063fa1434f4ae9bb63488cd79e2f76dea59e0e2d6cdec7236c2bb49ffb37da37cb7966be74eca5a171f659fee7bc501",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		fmt.Println(actualOutput)
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Fail retrieve deposit-data for unmatched index and publickey", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"deposit-data-seedless",
			"--private-key=63bc15d14d1460491535700fa2b6ac8873e1ede401cfc46e0c5ce77f00633d29",
			"--index=5",
			"--publickey=81fd26fe6e7cdbe1d0d45020050ba94c625f5236bf162b9ad3fca137d9120a0572c6f59b8cc70fae6cd6bb471b673e97",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "failed to get account by public key: account not found")
	})
}
