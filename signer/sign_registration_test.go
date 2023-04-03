package signer

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestSimpleSigner_SignRegistration(t *testing.T) {
	signer, err := setupNoSlashingProtectionSK(_byteArray("659e875e1b062c03f2f2a57332974d475b97df6cfc581d322e79642d39aca8fd"))
	require.NoError(t, err)

	valRegistrationMock := &apiv1.ValidatorRegistration{
		GasLimit:  123456,
		Timestamp: time.Unix(1658313712, 0),
	}
	copy(valRegistrationMock.FeeRecipient[:], _byteArray("9831EeF7A86C19E32bEcDad091c1DbC974cf452a"))
	copy(valRegistrationMock.Pubkey[:], _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"))

	registrationMock := &api.VersionedValidatorRegistration{
		Version: 0,
		V1:      valRegistrationMock,
	}

	tests := []struct {
		name          string
		data          *api.VersionedValidatorRegistration
		pubKey        []byte
		domain        [32]byte
		expectedError error
		sig           []byte
	}{
		{
			name:          "simple sign",
			data:          registrationMock,
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: nil,
			sig:           _byteArray("b088d9d27c783f3d5eb57a0df1e99f030e035ebcfdeb745da95400ab46a0c461f05f61533379d3bc56c5e94dfdf8560d0a31cfb9162f11ba9a82522f4043764a02008f6fef3b0167cbf2db9a749095343412a38568fe39d14c3ebcdddad7ee36"),
		},
		{
			name:          "nil data",
			data:          nil,
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("registration data is nil"),
			sig:           _byteArray("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
		},
		{
			name:          "unknown account, should error",
			data:          registrationMock,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
		{
			name:          "nil account, should error",
			data:          registrationMock,
			pubKey:        nil,
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("account was not supplied"),
			sig:           nil,
		},
		{
			name:          "empty account, should error",
			data:          registrationMock,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("00000001d7a9bca8823e555db65bb772e1496a26e1a8c5b1c0c7def9c9eaf7f6"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignRegistration(test.data, test.domain, test.pubKey)
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
