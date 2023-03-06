package signer

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// tested against a real block and sig from the Prater testnet (slot 5133683)
func TestSimpleSigner_SignVoluntaryExit(t *testing.T) {
	signer, err := setupNoSlashingProtectionSK(_byteArray("37247532b925101f094fb0cc877f523859c4a73bcbfc88f3833b05c26bd37cc6"))
	require.NoError(t, err)

	voluntaryExitMock := &phase0.VoluntaryExit{
		Epoch:          160427,
		ValidatorIndex: 438850,
	}

	tests := []struct {
		name          string
		data          *phase0.VoluntaryExit
		pubKey        []byte
		domain        [32]byte
		expectedError error
		sig           []byte
	}{
		{
			name:          "simple sign",
			data:          voluntaryExitMock,
			pubKey:        _byteArray("b5ade10d8cc63646ae7b30588c6fb9e482e51f98e396633a6e157bbde14bcdb771b7d147e5fb8b2bd6ce99323431008e"),
			domain:        _byteArray32("04000000c2ce3aa85707d491e3dd033a53971deb9bed9d4813d74c99369642f5"),
			expectedError: nil,
			sig:           _byteArray("8b8084ef095af3d0351a4c9308667b7254f3c0e9233e18f7ab59a29a6b6a3abdab9fbe9b7b61d9dd384675c4ed2b721a108890645ee9f69e97e2bccc586a35ddcebeaf20617d9c942fa1562db6814b016b8ebb4ee97d78c8ae27ae4b3dba2653"),
		},
		{
			name:          "nil data",
			data:          nil,
			pubKey:        _byteArray("b5ade10d8cc63646ae7b30588c6fb9e482e51f98e396633a6e157bbde14bcdb771b7d147e5fb8b2bd6ce99323431008e"),
			domain:        _byteArray32("04000000c2ce3aa85707d491e3dd033a53971deb9bed9d4813d74c99369642f5"),
			expectedError: errors.New("voluntary exit data is nil"),
			sig:           _byteArray("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
		},
		{
			name:          "unknown account, should error",
			data:          voluntaryExitMock,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("04000000c2ce3aa85707d491e3dd033a53971deb9bed9d4813d74c99369642f5"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
		{
			name:          "nil account, should error",
			data:          voluntaryExitMock,
			pubKey:        nil,
			domain:        _byteArray32("04000000c2ce3aa85707d491e3dd033a53971deb9bed9d4813d74c99369642f5"),
			expectedError: errors.New("account was not supplied"),
			sig:           nil,
		},
		{
			name:          "empty account, should error",
			data:          voluntaryExitMock,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("04000000c2ce3aa85707d491e3dd033a53971deb9bed9d4813d74c99369642f5"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignVoluntaryExit(test.data, test.domain, test.pubKey)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)
				require.EqualValues(t, test.sig, res)
			}
		})
	}
}
