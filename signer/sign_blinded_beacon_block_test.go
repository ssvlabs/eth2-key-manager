package signer

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electrta "github.com/attestantio/go-eth2-client/api/v1/electra"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/stretchr/testify/require"

	eth2keymanager "github.com/ssvlabs/eth2-key-manager"
	"github.com/ssvlabs/eth2-key-manager/core"
	prot "github.com/ssvlabs/eth2-key-manager/slashing_protection"
	"github.com/ssvlabs/eth2-key-manager/wallets"
)

type blindedBeaconBlockTest struct {
	name          string
	filename      string
	sk            string
	pk            string
	domain        string
	expectedSig   string
	unmarshalFunc func(data []byte) *api.VersionedBlindedBeaconBlock
}

func TestBlindedBeaconBlockProposals(t *testing.T) {
	require.NoError(t, core.InitBLS())

	tests := []blindedBeaconBlockTest{
		{
			name:        "bellatrix_blinded",
			sk:          "2550fa1883222bd14d84ff3794ed555661bc8d00e9305387b308ab4d285105b6",
			pk:          "99372f56eab6a56657006a831c6ffbad14a00b143eeb810364b0d6828f85414ae236980f93ef36a689f6d811352a6cf0",
			domain:      "00000000e7acb21061790987fa1c1e745cccfb358370b33e8af2b2c18938e6c2",
			expectedSig: "ae2605a0b64ababeaa41d27b60df1c716f9c180fb7b16d0d6c33743aa8ce3c387c870beec56dc7641d0824a130f3e68e1153cee3e243b033744092c833ac4fc86832b9158882f85747b1dfc0a9482924e7b40c4fd90fd65bcfe0f0cd75496971",
			unmarshalFunc: func(data []byte) *api.VersionedBlindedBeaconBlock {
				blk := &apiv1bellatrix.BlindedBeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &api.VersionedBlindedBeaconBlock{
					Version:   spec.DataVersionBellatrix,
					Bellatrix: blk,
				}
			},
		},
		{
			name:        "capella_blinded",
			sk:          "17f8b2dbb824318970a9c82b4f5456d598342dac248f89c252dd30bb9e85e43c",
			pk:          "81d6fc2f01633e8eab3ba4d72588e14f45b00e68ab887bdd4ec5e8558965db21189310df973837106216777b07fc0805",
			domain:      "0000000047eb72b3be36f08feffcaba760f0a2ed78c1a85f0654941a0d19d0fa",
			expectedSig: "9321eb1d373e75bb0b746a8c4b6f62285cccb17f2105a1f4f968534902131a29fffd5b54a0012c34c8600639609c46ab05a23fe8f89cdd252df680883ccb94073d9b852680bba6c5f5166edcd6024c9ef489c555d1cde994fdd7010fb394e45d",
			unmarshalFunc: func(data []byte) *api.VersionedBlindedBeaconBlock {
				blk := &apiv1capella.BlindedBeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &api.VersionedBlindedBeaconBlock{
					Version: spec.DataVersionCapella,
					Capella: blk,
				}
			},
		},
		{
			name:        "deneb_blinded",
			sk:          "1e27ca2fc62c9904662dbc736a28146d570d2289443f6ffec03917e557eaa179",
			pk:          "8b3bb838e8adeff9dbf11bd63d3767222b2b5f34e749fd5599bb093afc040bc86b6950243c4d41d2aa7526534489e355",
			domain:      "00000000a75dccf2e267634e64c6cb6df009db33c9b870393d69264778a9a195",
			expectedSig: "b10b2689062b3bf14333aada3a5abe56cdf00f665f9e0c40cf6c20ccda62067c1f3a7f4d2dc749966f83d166fb935e8413264138b697fcd33036efe315bc7ef4aae98e9898ad4c92507ecc4b10ba85afb622778e4b2dbb4c4bc9f173573703a4",
			unmarshalFunc: func(data []byte) *api.VersionedBlindedBeaconBlock {
				blk := &apiv1deneb.BlindedBeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &api.VersionedBlindedBeaconBlock{
					Version: spec.DataVersionDeneb,
					Deneb:   blk,
				}
			},
		},
		{
			name:        "electra_blinded",
			sk:          "4387296cbebc40a5dbb3a5b39b65a0c58a7e673944d44da3602e061762dd028d",
			pk:          "b846c9e4e4482aece0dd3f82b3ac429e9315542953f2960f1575f239d3efccc73fa783591378cf349f575369932a2b8b",
			domain:      "00000000e551da7785705c8c2e2f4ed229e22096eab18686d02a647950a11205",
			expectedSig: "adb8c16ce9f9505e7faa495fcda269c4f9a6498e61c4b716b9afb347f60593ee075c315c7f65dd677ea226bd40da2d430f6948b8b255c288577284dc0d3bdecc50cf8ab46044b3e4ad102b5d351555133ff86968f3b1ec4e3f6fac331b8cdee3",
			unmarshalFunc: func(data []byte) *api.VersionedBlindedBeaconBlock {
				blk := &apiv1electrta.BlindedBeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &api.VersionedBlindedBeaconBlock{
					Version: spec.DataVersionElectra,
					Electra: blk,
				}
			},
		},
		{
			name:        "fulu_blinded",
			filename:    "electra_blinded", // reuse electra block
			sk:          "4387296cbebc40a5dbb3a5b39b65a0c58a7e673944d44da3602e061762dd028d",
			pk:          "b846c9e4e4482aece0dd3f82b3ac429e9315542953f2960f1575f239d3efccc73fa783591378cf349f575369932a2b8b",
			domain:      "00000000e551da7785705c8c2e2f4ed229e22096eab18686d02a647950a11205",
			expectedSig: "adb8c16ce9f9505e7faa495fcda269c4f9a6498e61c4b716b9afb347f60593ee075c315c7f65dd677ea226bd40da2d430f6948b8b255c288577284dc0d3bdecc50cf8ab46044b3e4ad102b5d351555133ff86968f3b1ec4e3f6fac331b8cdee3",
			unmarshalFunc: func(data []byte) *api.VersionedBlindedBeaconBlock {
				blk := &apiv1electrta.BlindedBeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &api.VersionedBlindedBeaconBlock{
					Version: spec.DataVersionFulu,
					Fulu:    blk,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.filename
			if filename == "" {
				filename = tt.name
			}
			regularFile := fmt.Sprintf("testdata/blocks/%s.json", filename)

			blockJSON, err := os.ReadFile(regularFile)
			require.NoError(t, err, "Failed to read JSON file for %s", tt.name)

			// Setup KeyVault
			store := inmemStorage()
			options := &eth2keymanager.KeyVaultOptions{}
			options.SetStorage(store)
			options.SetWalletType(core.NDWallet)
			vault, err := eth2keymanager.NewKeyVault(options)
			require.NoError(t, err)

			wallet, err := vault.Wallet()
			require.NoError(t, err)

			// Create account
			k, err := core.NewHDKeyFromPrivateKey(_byteArray(tt.sk), "")
			require.NoError(t, err)
			acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
			require.NoError(t, wallet.AddValidatorAccount(acc))
			require.Equal(t, _byteArray(tt.pk), acc.ValidatorPublicKey())

			// Setup signer
			signer := NewSimpleSigner(wallet, &prot.NoProtection{}, core.PraterNetwork)

			block := tt.unmarshalFunc(blockJSON)
			require.NoError(t, err, "Failed to unmarshal block for %s", tt.name)
			sig, _, err := signer.SignBlindedBeaconBlock(block, _byteArray32(tt.domain), _byteArray(tt.pk))
			require.NoError(t, err, "Failed to sign block for %s", tt.name)
			require.EqualValues(t, _byteArray(tt.expectedSig), sig, "Signature mismatch for %s", tt.name)
		})
	}
}

// Test slashing by signing first beacon block and then blinded beacon block
func TestDoubleProposalsSigning_Regular_Blinded(t *testing.T) {
	require.NoError(t, core.InitBLS())

	sk := "2550fa1883222bd14d84ff3794ed555661bc8d00e9305387b308ab4d285105b6"
	pk := "99372f56eab6a56657006a831c6ffbad14a00b143eeb810364b0d6828f85414ae236980f93ef36a689f6d811352a6cf0"
	domain := "00000000e7acb21061790987fa1c1e745cccfb358370b33e8af2b2c18938e6c2"
	sigByts := "ae2605a0b64ababeaa41d27b60df1c716f9c180fb7b16d0d6c33743aa8ce3c387c870beec56dc7641d0824a130f3e68e1153cee3e243b033744092c833ac4fc86832b9158882f85747b1dfc0a9482924e7b40c4fd90fd65bcfe0f0cd75496971"

	blockJSON, err := os.ReadFile("testdata/blocks/bellatrix.json")
	require.NoError(t, err)

	blindedBlockJSON, err := os.ReadFile("testdata/blocks/bellatrix_blinded.json")
	require.NoError(t, err)

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	protector := prot.NewNormalProtection(store)

	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(_byteArray(sk), "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	signer := NewSimpleSigner(wallet, protector, core.PraterNetwork)

	// decode block
	blk := &bellatrix.BeaconBlock{}
	require.NoError(t, blk.UnmarshalJSON(blockJSON))
	// minimal slashing protection
	require.NoError(t, protector.UpdateHighestProposal(_byteArray(pk), blk.Slot-1))

	versionedBeaconBlock := &spec.VersionedBeaconBlock{
		Version:   spec.DataVersionBellatrix,
		Bellatrix: blk,
	}

	// decode blinded block
	blindedBlk := &apiv1bellatrix.BlindedBeaconBlock{}
	require.NoError(t, blindedBlk.UnmarshalJSON(blindedBlockJSON))
	versionedBlindedBeaconBlock := &api.VersionedBlindedBeaconBlock{
		Version:   spec.DataVersionBellatrix,
		Bellatrix: blindedBlk,
	}

	sig, _, err := signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32(domain), _byteArray(pk))
	require.NoError(t, err)
	require.EqualValues(t, _byteArray(sigByts), sig)

	_, _, err = signer.SignBlindedBeaconBlock(versionedBlindedBeaconBlock, _byteArray32(domain), _byteArray(pk))
	require.Error(t, err)
	require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
}

// Test slashing by signing first blinded beacon block and then beacon block
func TestDoubleProposalsSigning_Blinded_Regular(t *testing.T) {
	require.NoError(t, core.InitBLS())

	sk := "2550fa1883222bd14d84ff3794ed555661bc8d00e9305387b308ab4d285105b6"
	pk := "99372f56eab6a56657006a831c6ffbad14a00b143eeb810364b0d6828f85414ae236980f93ef36a689f6d811352a6cf0"
	domain := "00000000e7acb21061790987fa1c1e745cccfb358370b33e8af2b2c18938e6c2"
	sigByts := "ae2605a0b64ababeaa41d27b60df1c716f9c180fb7b16d0d6c33743aa8ce3c387c870beec56dc7641d0824a130f3e68e1153cee3e243b033744092c833ac4fc86832b9158882f85747b1dfc0a9482924e7b40c4fd90fd65bcfe0f0cd75496971"

	blockJSON, err := os.ReadFile("testdata/blocks/bellatrix.json")
	require.NoError(t, err)

	blindedBlockJSON, err := os.ReadFile("testdata/blocks/bellatrix_blinded.json")
	require.NoError(t, err)

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	protector := prot.NewNormalProtection(store)

	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(_byteArray(sk), "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	signer := NewSimpleSigner(wallet, protector, core.PraterNetwork)

	// decode block
	blk := &bellatrix.BeaconBlock{}
	require.NoError(t, blk.UnmarshalJSON(blockJSON))

	// minimal slashing protection
	require.NoError(t, protector.UpdateHighestProposal(_byteArray(pk), blk.Slot-1))

	versionedBeaconBlock := &spec.VersionedBeaconBlock{
		Version:   spec.DataVersionBellatrix,
		Bellatrix: blk,
	}

	// decode blinded block
	blindedBlk := &apiv1bellatrix.BlindedBeaconBlock{}
	require.NoError(t, blindedBlk.UnmarshalJSON(blindedBlockJSON))
	versionedBlindedBeaconBlock := &api.VersionedBlindedBeaconBlock{
		Version:   spec.DataVersionBellatrix,
		Bellatrix: blindedBlk,
	}

	sig, _, err := signer.SignBlindedBeaconBlock(versionedBlindedBeaconBlock, _byteArray32(domain), _byteArray(pk))
	require.NoError(t, err)
	require.EqualValues(t, _byteArray(sigByts), sig)

	_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32(domain), _byteArray(pk))
	require.Error(t, err)
	require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
}
