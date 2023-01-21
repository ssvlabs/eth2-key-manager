package signer

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	prot "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/bloxapp/eth2-key-manager/wallets"
)

func TestReferenceAttestationAggregation(t *testing.T) {
	sk := _byteArray("2c083f2c8fc923fa2bd32a70ab72b4b46247e8c1f347adc30b2f8036a355086c")
	pk := _byteArray("a9cf360aa15fb1d1d30ee2b578dc5884823c19661886ae8b892775ccb3bd96b7d7345569a2aa0b14e4d015c54a6a0c54")
	aggByts := _byteArray("01000000000000006c000000b4fa352d2d6dbdf884266af7ea0914451929b343527ea6c1737ac93b3dde8b7c98e6ce61d68b7a2e7b7af8f8d0fd429d0bdd5f930b83e6842bf4342d3d1d3d10fc0d15bab7649bb8aa8287ca104a1f79d396ce0217bb5cd3e6503a3bce4c9776e4000000000000000000000000000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776b0000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000003a43a4bf26fb5947e809c1f24f7dc6857c8ac007e535d48e6e4eca2122fd776bb4fa352d2d6dbdf884266af7ea0914451929b343527ea6c1737ac93b3dde8b7c98e6ce61d68b7a2e7b7af8f8d0fd429d0bdd5f930b83e6842bf4342d3d1d3d10fc0d15bab7649bb8aa8287ca104a1f79d396ce0217bb5cd3e6503a3bce4c97760010")
	domain := _byteArray32("0100000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459")
	sig := _byteArray("a8333dee0d7a20d99d13f870c3b77e413b7755c1640985fc70bc58f6004b50f43ef301147208c9c5393258d7e9b2208316c48d540879e2818352b346d8ce6d91ce8c1942758ab5425a0448959ea46609397c0dc9c05708f243389af694fda91c")

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(sk, "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, err)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	signer := NewSimpleSigner(wallet, &prot.NoProtection{}, core.PraterNetwork)

	// decode attestation
	agg := &phase0.AggregateAndProof{}
	require.NoError(t, agg.UnmarshalSSZ(aggByts))

	actualSig, err := signer.SignAggregateAndProof(agg, domain, pk)
	require.NoError(t, err)
	require.EqualValues(t, sig, actualSig)
}
