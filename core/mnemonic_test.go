package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

func TestSeedFromMnemonic(t *testing.T) {
	e2types.InitBLS()

	tests := []struct {
		mnemonic        string
		password        string
		expectedSeedHex string
	}{
		{
			mnemonic:        "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless",
			password:        "TREZOR",
			expectedSeedHex: "c0c519bd0e91a2ed54357d9d1ebef6f5af218a153624cf4f2da911a0ed8f7a09e2ef61af0aca007096df430022f7a2b6fb91661a9589097069720d015e4e982f",
		},
		{
			mnemonic:        "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo vote",
			password:        "TREZOR",
			expectedSeedHex: "dd48c104698c30cfe2b6142103248622fb7bb0ff692eebb00089b32d22484e1613912f0a5b694407be899ffd31ed3992c456cdf60f5d4564b8ba3f05a69890ad",
		},
		{
			mnemonic:        "gravity machine north sort system female filter attitude volume fold club stay feature office ecology stable narrow fog",
			password:        "TREZOR",
			expectedSeedHex: "628c3827a8823298ee685db84f55caa34b5cc195a778e52d45f59bcf75aba68e4d7590e101dc414bc1bbd5737666fbbef35d1f1903953b66624f910feef245ac",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test vector: %d", i), func(t *testing.T) {
			expectedSeed, err := hex.DecodeString(test.expectedSeedHex)
			require.NoError(t, err)

			// KeyVault
			fromMnemonic, err := SeedFromMnemonic(test.mnemonic, test.password)
			require.NoError(t, err)

			require.Equal(t, expectedSeed, fromMnemonic)
		})
	}
}

func TestGenerateNewEntropy(t *testing.T) {
	got, err := GenerateNewEntropy()
	require.NoError(t, err)
	require.NotEmpty(t, got)
}

func TestEntropyToMnemonic(t *testing.T) {
	t.Run("successfully generated mnemonic", func(t *testing.T) {
		entropy, err := GenerateNewEntropy()
		require.NoError(t, err)
		require.NotEmpty(t, entropy)

		mnemonic, err := EntropyToMnemonic(entropy)
		require.NoError(t, err)
		require.NotEmpty(t, mnemonic)
	})

	t.Run("rejects generate mnemonic with empty entropy", func(t *testing.T) {
		mnemonic, err := EntropyToMnemonic(nil)
		require.Error(t, err)
		require.Empty(t, mnemonic)
	})
}
