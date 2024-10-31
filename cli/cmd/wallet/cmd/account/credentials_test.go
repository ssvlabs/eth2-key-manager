package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/util/printer"
)

func TestAccountCredentials(t *testing.T) {
	t.Run("Successfully handle credentials change at specific index", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--validator-indices=273230",
			"--validator-public-keys=0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Successfully handle accumulated credentials change", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-indices=273230,273407",
			"--validator-public-keys=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c,0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d15,0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e,0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)
	})

	t.Run("Only one validator can be specified if accumulate is false", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=false",
			"--validator-indices=273230,273407",
			"--validator-public-keys=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c,0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d15,0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e,0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect credentials flags: only one validator can be specified if accumulate is false")
	})

	t.Run("Not equal length - should be 2 validator indices", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-indices=273230,273407,273407",
			"--validator-public-keys=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c,0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d15,0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e,0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect credentials flags: validator indices, public keys, withdrawal credentials and to execution addresses must be of equal length")
	})

	t.Run("Not equal length - should be two public keys", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-indices=273230,273407",
			"--validator-public-keys=0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d15,0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e,0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect credentials flags: validator indices, public keys, withdrawal credentials and to execution addresses must be of equal length")
	})

	t.Run("Not equal length - should be two withdrawal credentials", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-indices=273230,273407",
			"--validator-public-keys=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c,0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e,0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect credentials flags: validator indices, public keys, withdrawal credentials and to execution addresses must be of equal length")
	})

	t.Run("Not equal length - should be two to execution addresses", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-indices=273230,273407",
			"--validator-public-keys=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c,0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d15,0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "failed to collect credentials flags: validator indices, public keys, withdrawal credentials and to execution addresses must be of equal length")
	})

	t.Run("Derived pub key does not match with the provided", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-indices=273230,273407",
			"--validator-public-keys=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7a,0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d15,0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e,0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "derived validator public key: 0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c, does not match with the provided one: 0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7a")
	})

	t.Run("Derived withdrawal credentials does not match with the provided", func(t *testing.T) {
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs([]string{
			"wallet",
			"account",
			"credentials",
			"--seed=847d135b3aecac8ae77c3fdfd46dc5849ad3b5bacd30a1b9082b6ff53c77357e923b12fcdc3d02728fd35c3685de1fe1e9c052c48f0d83566b1b2287cf0e54c3",
			"--index=1",
			"--accumulate=true",
			"--validator-indices=273230,273407",
			"--validator-public-keys=0xb2dc1daa8c9cd104d4503028639e41a41e4f06ee5cc90ebfaeab3c41f43a148ce9afa4ebd1b8be3f54e4d6c15e870c7c,0xa1a593775967bf88bb6c14ac109c12e52dc836fa139bd1ba6ca873d65fe91bb7a0fc79c7b2a7315482a81f31e6b1018a",
			"--withdrawal-credentials=0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d14,0x00202f88c6116d27f5de06eeda3b801d3a4eeab8eb09f879848445fffc8948a5",
			"--to-execution-address=0x3e6935b8250Cf9A777862871649E5594bE08779e,0x3e6935b8250Cf9A777862871649E5594bE08779e",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.EqualValues(t, actualOutput, "")
		require.Error(t, err)
		require.EqualError(t, err, "derived withdrawal credentials: 0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d15, does not match with the provided one: 0x00d9cdf17e3a79317a4e5cd18580b1d10b1df360bbca5c5f8ac5b79b45c29d14")
	})
}
