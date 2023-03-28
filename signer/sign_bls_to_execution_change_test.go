package signer

import (
	"encoding/hex"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// tested against a real block and sig from the Sepolia testnet (slot 1864649)
func TestSimpleSigner_SignBLSToExecutionChange(t *testing.T) {
	signer, err := setupNoSlashingProtectionSK(_byteArray("6c1efcf889f78ec02d6becac5839b656d402d68f1c56723616f2c22d69cb7fc1"))
	require.NoError(t, err)

	blsToExecutionChangeMock := &capella.BLSToExecutionChange{
		ValidatorIndex: 1573,
	}
	hexDecodedExecAdd, err := hex.DecodeString("beefd32838d5762ff55395a7beebef6e8528c64f")
	require.NoError(t, err)
	hexDecodedBLSPubKey, err := hex.DecodeString("8d8e66062fa5a1e5c4b9b0017d4027c944550a4e096fd8c535de1aa3b0283c5ece23c68c5881a154fe24a9c9377a0a09")
	require.NoError(t, err)

	copy(blsToExecutionChangeMock.ToExecutionAddress[:], hexDecodedExecAdd)
	copy(blsToExecutionChangeMock.FromBLSPubkey[:], hexDecodedBLSPubKey)

	tests := []struct {
		name          string
		data          *capella.BLSToExecutionChange
		pubKey        []byte
		domain        [32]byte
		expectedError error
		sig           []byte
	}{
		{
			name:          "simple sign",
			data:          blsToExecutionChangeMock,
			pubKey:        blsToExecutionChangeMock.FromBLSPubkey[:],
			domain:        _byteArray32("0a000000a8fee8ee9978418b64f1140b05f699a49ccd9b3fd666c35d4ae5f79e"),
			expectedError: nil,
			sig:           _byteArray("aae6b0261494230fbf69ec8c1b907763153a5c52a39797f90aa106923d2d0d4752a392642f555520c2bbf54a9191876e0642555afa1c0050b341314c610f5bfd0821eafcc7981d551cc5b4969aff0ede5c229084dd098244706336621809a069"),
		},
		{
			name:          "nil data",
			data:          nil,
			pubKey:        blsToExecutionChangeMock.FromBLSPubkey[:],
			domain:        _byteArray32("0a000000a8fee8ee9978418b64f1140b05f699a49ccd9b3fd666c35d4ae5f79e"),
			expectedError: errors.New("bls to execution change is nil"),
			sig:           _byteArray("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
		},
		{
			name:          "unknown account, should error",
			data:          blsToExecutionChangeMock,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0a000000a8fee8ee9978418b64f1140b05f699a49ccd9b3fd666c35d4ae5f79e"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
		{
			name:          "nil account, should error",
			data:          blsToExecutionChangeMock,
			pubKey:        nil,
			domain:        _byteArray32("0a000000a8fee8ee9978418b64f1140b05f699a49ccd9b3fd666c35d4ae5f79e"),
			expectedError: errors.New("account was not supplied"),
			sig:           nil,
		},
		{
			name:          "empty account, should error",
			data:          blsToExecutionChangeMock,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignBLSToExecutionChange(test.data, test.domain, test.pubKey)
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
