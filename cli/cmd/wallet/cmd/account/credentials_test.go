package account_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
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

	// The signatures in this test were verified with Prysm's libraries, using the following code:
	//
	//  func TestVerifyBLSChangeSignature(t *testing.T) {
	//  	params.SetActive(params.MainnetConfig())
	//  	change := &ethpb.SignedBLSToExecutionChange{
	//  		Message: &ethpb.BLSToExecutionChange{
	//  			ValidatorIndex:     273230,
	//  			FromBlsPubkey:      hexutil.MustDecode("0x8c02c584a9265f6ff2a2119c4eaee0385fb4320aa6ddb14ad98050af7bc3e250aadafd8d4733f980408e8a595583db6e"),
	//  			ToExecutionAddress: hexutil.MustDecode("0x3e6935b8250cf9a777862871649e5594be08779e"),
	//  		},
	//  		Signature: hexutil.MustDecode("0xb0bd98cfe89b8439a8e787e25fb25f6999db397bfd47d34eca9e4a85bca6b0c529fa83fc8705fcbba95ce7623713a60610a7b0ba14c6816fe74f20f00c3226192a064529047ae5e4009cfa62273e9f9de5d03806ff6019443c10fc10a7fb0d11"),
	//  	}
	//  	spb := &ethpb.BeaconStateCapella{
	//  		Fork: &ethpb.Fork{
	//  			CurrentVersion:  params.BeaconConfig().CapellaForkVersion,
	//  			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
	//  			Epoch:           params.BeaconConfig().CapellaForkEpoch,
	//  		},
	//  		GenesisValidatorsRoot: hexutil.MustDecode("0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95"),
	//  	}
	//  	st, err := state_native.InitializeFromProtoCapella(spb)
	//  	require.NoError(t, err)
	//  	changeV2 := migration.V1Alpha1SignedBLSToExecChangeToV2(change)
	//  	require.NoError(t, blocks.VerifyBLSChangeSignature(st, changeV2))
	//  }
	t.Run("Successfully handle accumulated credentials change (mainnet)", func(t *testing.T) {
		input := []string{
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
			"--network=mainnet",
		}
		expectedOutput := `[
			{
			  "message": {
				"validator_index": "273230",
				"from_bls_pubkey": "0x8c02c584a9265f6ff2a2119c4eaee0385fb4320aa6ddb14ad98050af7bc3e250aadafd8d4733f980408e8a595583db6e",
				"to_execution_address": "0x3e6935b8250cf9a777862871649e5594be08779e"
			  },
			  "signature": "0xb0bd98cfe89b8439a8e787e25fb25f6999db397bfd47d34eca9e4a85bca6b0c529fa83fc8705fcbba95ce7623713a60610a7b0ba14c6816fe74f20f00c3226192a064529047ae5e4009cfa62273e9f9de5d03806ff6019443c10fc10a7fb0d11"
			},
			{
			  "message": {
				"validator_index": "273407",
				"from_bls_pubkey": "0x862418bfdd18e1147e6fb62c9c7cddc638cc74819d7fbecb947a571145e7782bbb44f8ef9fa833410ef4e9dae2903756",
				"to_execution_address": "0x3e6935b8250cf9a777862871649e5594be08779e"
			  },
			  "signature": "0x8547adbba39a4c600c8308109ddc1a4705b34ba3b419fd52570b00a1d0ef811149620c7a7e85390f8a920d1458773b5b0cfc9b52d7c2421e0cd36119e6ce156276c906504b4ba4a8431e9fe717b30484666f9607cacdcf8eaa1830902c678721"
			}
		]`

		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs(input)
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)

		require.JSONEq(t, expectedOutput, actualOutput)
	})

	t.Run("Successfully handle accumulated credentials change (prater)", func(t *testing.T) {
		input := []string{
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
		}
		expectedOutput := `[
			{
			  "message": {
				"validator_index": "273230",
				"from_bls_pubkey": "0x8c02c584a9265f6ff2a2119c4eaee0385fb4320aa6ddb14ad98050af7bc3e250aadafd8d4733f980408e8a595583db6e",
				"to_execution_address": "0x3e6935b8250cf9a777862871649e5594be08779e"
			  },
			  "signature": "0x8265248c611cfd8202d82c42a222466051d034a01e5bd48c3c29ff47c7e557f5cbcb32713cf205ca61bcad5f59ea7df4018ff863d4318b53ea3b2383737149c16a24bea61784de5deb5ba9f5ea78f4ef7950b47e5aafa6805f238d20b80b7de7"
			},
			{
			  "message": {
				"validator_index": "273407",
				"from_bls_pubkey": "0x862418bfdd18e1147e6fb62c9c7cddc638cc74819d7fbecb947a571145e7782bbb44f8ef9fa833410ef4e9dae2903756",
				"to_execution_address": "0x3e6935b8250cf9a777862871649e5594be08779e"
			  },
			  "signature": "0x969693dc21fd9b05fb9737d16909ec1c9b3f5c8b29f1d4c75b1e96609fa20b926ece5a9d8a1be4a2601bac0d514b523d0affd2f8e576c42718f08ddbcb806dfcb6becf832f1cb31dd8e474a0d9d5838564f43a64fcf54d46b409aa94de32195f"
			}
		]`

		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		cmd.RootCmd.SetArgs(input)
		err := cmd.RootCmd.Execute()
		actualOutput := output.String()
		require.NotNil(t, actualOutput)
		require.NoError(t, err)

		require.JSONEq(t, expectedOutput, actualOutput)
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
