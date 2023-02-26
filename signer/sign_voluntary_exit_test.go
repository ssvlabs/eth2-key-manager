package signer

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestSimpleSigner_SignVoluntaryExit(t *testing.T) {
	signer, err := setupNoSlashingProtectionSK(_byteArray("659e875e1b062c03f2f2a57332974d475b97df6cfc581d322e79642d39aca8fd"))
	require.NoError(t, err)

	voluntaryExitMock := &phase0.VoluntaryExit{
		Epoch:          1,
		ValidatorIndex: 0,
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
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: nil,
			sig:           _byteArray("895740a6edec2907d16cc53b8c1357f8984706553a470748df2577cc6d881c6b75f88337bfad30421f9d620bb1dcb4ce15efa29dfc38679e2b3d3e99e0d773421ccb67f650522af1dac606327b2dacce8e5d767c6e4a6ed1eca45170d0a07c3c"),
		},
		{
			name:          "nil data",
			data:          nil,
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("voluntary exit data is nil"),
			sig:           _byteArray("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
		},
		{
			name:          "unknown account, should error",
			data:          voluntaryExitMock,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
		{
			name:          "nil account, should error",
			data:          voluntaryExitMock,
			pubKey:        nil,
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("account was not supplied"),
			sig:           nil,
		},
		{
			name:          "empty account, should error",
			data:          voluntaryExitMock,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignVoluntaryExit(test.data, test.domain, test.pubKey)
			fmt.Println(hex.EncodeToString(res))
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
