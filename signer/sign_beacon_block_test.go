package signer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	eth2keymanager "github.com/ssvlabs/eth2-key-manager"
	"github.com/ssvlabs/eth2-key-manager/core"
	prot "github.com/ssvlabs/eth2-key-manager/slashing_protection"
	"github.com/ssvlabs/eth2-key-manager/wallets"
)

func testBlock(t *testing.T) *phase0.BeaconBlock {
	blockJSON, err := os.ReadFile("testdata/blocks/phase0.json")
	require.NoError(t, err)

	blk := &phase0.BeaconBlock{}
	require.NoError(t, json.Unmarshal(blockJSON, blk))
	return blk
}

type beaconBlockTest struct {
	name          string
	filename      string
	sk            string
	pk            string
	domain        string
	expectedSig   string
	unmarshalFunc func(data []byte) *spec.VersionedBeaconBlock
}

func TestBeaconBlockProposals(t *testing.T) {
	require.NoError(t, core.InitBLS())

	tests := []beaconBlockTest{
		{
			name:        "phase0",
			sk:          "5470813f7deef638dc531188ca89e36976d536f680e89849cd9077fd096e20bc",
			pk:          "a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972",
			domain:      "0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459",
			expectedSig: "aab8d1dc7e9fdea39b9a46cde64759138372a8d1226d56d01fa9e2df9e45e39d67bacf5794bb6841f5ea143aa124a33211089d4cc045e0605380bf4ae03291b1ee8082e55274ce5dc513cc0bebdb8e8a5276ad39a0cea3edd321d9cf854f695c",
			unmarshalFunc: func(data []byte) *spec.VersionedBeaconBlock {
				blk := &phase0.BeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &spec.VersionedBeaconBlock{
					Version: spec.DataVersionPhase0,
					Phase0:  blk,
				}
			},
		},
		{
			name:        "altair",
			sk:          "2799ceccbdaf1e36679b413193a363bfe6d2d35c8cf6ff6151165707461eaed7",
			pk:          "b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d8",
			domain:      "000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e",
			expectedSig: "ae37c5b026490df77d4757ad278f1a5f2c46c8971eaa9b62ad09534d3626b1d1da1371edeb061ed6f2f327f79d3bf79101941225dfb97af6d15c8fd6ac7ae8a180392a19aecec50c938b1ea6faf8ba848a4dc9a01f46502e98caec4a01b06575",
			unmarshalFunc: func(data []byte) *spec.VersionedBeaconBlock {
				blk := &altair.BeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &spec.VersionedBeaconBlock{
					Version: spec.DataVersionAltair,
					Altair:  blk,
				}
			},
		},
		{
			name:        "bellatrix",
			sk:          "2550fa1883222bd14d84ff3794ed555661bc8d00e9305387b308ab4d285105b6",
			pk:          "99372f56eab6a56657006a831c6ffbad14a00b143eeb810364b0d6828f85414ae236980f93ef36a689f6d811352a6cf0",
			domain:      "00000000e7acb21061790987fa1c1e745cccfb358370b33e8af2b2c18938e6c2",
			expectedSig: "ae2605a0b64ababeaa41d27b60df1c716f9c180fb7b16d0d6c33743aa8ce3c387c870beec56dc7641d0824a130f3e68e1153cee3e243b033744092c833ac4fc86832b9158882f85747b1dfc0a9482924e7b40c4fd90fd65bcfe0f0cd75496971",
			unmarshalFunc: func(data []byte) *spec.VersionedBeaconBlock {
				blk := &bellatrix.BeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &spec.VersionedBeaconBlock{
					Version:   spec.DataVersionBellatrix,
					Bellatrix: blk,
				}
			},
		},
		{
			name:        "capella",
			sk:          "17f8b2dbb824318970a9c82b4f5456d598342dac248f89c252dd30bb9e85e43c",
			pk:          "81d6fc2f01633e8eab3ba4d72588e14f45b00e68ab887bdd4ec5e8558965db21189310df973837106216777b07fc0805",
			domain:      "0000000047eb72b3be36f08feffcaba760f0a2ed78c1a85f0654941a0d19d0fa",
			expectedSig: "9321eb1d373e75bb0b746a8c4b6f62285cccb17f2105a1f4f968534902131a29fffd5b54a0012c34c8600639609c46ab05a23fe8f89cdd252df680883ccb94073d9b852680bba6c5f5166edcd6024c9ef489c555d1cde994fdd7010fb394e45d",
			unmarshalFunc: func(data []byte) *spec.VersionedBeaconBlock {
				blk := &capella.BeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &spec.VersionedBeaconBlock{
					Version: spec.DataVersionCapella,
					Capella: blk,
				}
			},
		},
		{
			name:        "deneb",
			sk:          "1e27ca2fc62c9904662dbc736a28146d570d2289443f6ffec03917e557eaa179",
			pk:          "8b3bb838e8adeff9dbf11bd63d3767222b2b5f34e749fd5599bb093afc040bc86b6950243c4d41d2aa7526534489e355",
			domain:      "00000000a75dccf2e267634e64c6cb6df009db33c9b870393d69264778a9a195",
			expectedSig: "b10b2689062b3bf14333aada3a5abe56cdf00f665f9e0c40cf6c20ccda62067c1f3a7f4d2dc749966f83d166fb935e8413264138b697fcd33036efe315bc7ef4aae98e9898ad4c92507ecc4b10ba85afb622778e4b2dbb4c4bc9f173573703a4",
			unmarshalFunc: func(data []byte) *spec.VersionedBeaconBlock {
				blk := &deneb.BeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &spec.VersionedBeaconBlock{
					Version: spec.DataVersionDeneb,
					Deneb:   blk,
				}
			},
		},
		{
			name:        "electra",
			sk:          "4387296cbebc40a5dbb3a5b39b65a0c58a7e673944d44da3602e061762dd028d",
			pk:          "b846c9e4e4482aece0dd3f82b3ac429e9315542953f2960f1575f239d3efccc73fa783591378cf349f575369932a2b8b",
			domain:      "00000000e551da7785705c8c2e2f4ed229e22096eab18686d02a647950a11205",
			expectedSig: "adb8c16ce9f9505e7faa495fcda269c4f9a6498e61c4b716b9afb347f60593ee075c315c7f65dd677ea226bd40da2d430f6948b8b255c288577284dc0d3bdecc50cf8ab46044b3e4ad102b5d351555133ff86968f3b1ec4e3f6fac331b8cdee3",
			unmarshalFunc: func(data []byte) *spec.VersionedBeaconBlock {
				blk := &electra.BeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &spec.VersionedBeaconBlock{
					Version: spec.DataVersionElectra,
					Electra: blk,
				}
			},
		},
		{
			name:        "fulu",
			filename:    "electra", // reuse electra block
			sk:          "4387296cbebc40a5dbb3a5b39b65a0c58a7e673944d44da3602e061762dd028d",
			pk:          "b846c9e4e4482aece0dd3f82b3ac429e9315542953f2960f1575f239d3efccc73fa783591378cf349f575369932a2b8b",
			domain:      "00000000e551da7785705c8c2e2f4ed229e22096eab18686d02a647950a11205",
			expectedSig: "adb8c16ce9f9505e7faa495fcda269c4f9a6498e61c4b716b9afb347f60593ee075c315c7f65dd677ea226bd40da2d430f6948b8b255c288577284dc0d3bdecc50cf8ab46044b3e4ad102b5d351555133ff86968f3b1ec4e3f6fac331b8cdee3",
			unmarshalFunc: func(data []byte) *spec.VersionedBeaconBlock {
				blk := &electra.BeaconBlock{}
				require.NoError(t, json.Unmarshal(data, blk))
				return &spec.VersionedBeaconBlock{
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
			sig, _, err := signer.SignBeaconBlock(block, _byteArray32(tt.domain), _byteArray(tt.pk))
			require.NoError(t, err, "Failed to sign block for %s", tt.name)
			require.EqualValues(t, _byteArray(tt.expectedSig), sig, "Signature mismatch for %s", tt.name)
		})
	}
}

func TestProposalSlashingSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(t, seed, true, true)
	require.NoError(t, err)

	t.Run("valid proposal", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}
		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})

	t.Run("valid proposal, sign using nil pk. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}
		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), nil)
		require.NotNil(t, err)
		require.Error(t, err, "account was not supplied")
	})

	t.Run("double proposal, different state root. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		blk.StateRoot = _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459")
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}

		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})

	t.Run("double proposal, different body root. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}

		copy(blk.Body.Graffiti[:], "different body root")
		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("domain"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})

	t.Run("double proposal, different parent root. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		blk.ParentRoot = _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52458")
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}

		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("domain"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})

	t.Run("double proposal, different proposer index. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		blk.ProposerIndex = 3
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}

		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("domain"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})
}

func TestFarFutureProposalSignature(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	network := core.PraterNetwork
	maxValidSlot := network.EstimatedSlotAtTime(time.Now().Add(MaxFarFutureDelta))

	t.Run("max valid source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		blk := testBlock(t)
		blk.Slot = maxValidSlot
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}

		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})
	t.Run("too far into the future source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(t, seed, true, true)
		require.NoError(t, err)

		blk := testBlock(t)
		blk.Slot = maxValidSlot + 1
		versionedBeaconBlock := &spec.VersionedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0:  blk,
		}

		_, _, err = signer.SignBeaconBlock(versionedBeaconBlock, _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.EqualError(t, err, "proposed block slot too far into the future")
	})
}
