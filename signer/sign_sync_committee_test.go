package signer

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestSimpleSigner_SignSyncCommitteeContributionAndProof(t *testing.T) {
	signer, err := setupNoSlashingProtectionSK(_byteArray("659e875e1b062c03f2f2a57332974d475b97df6cfc581d322e79642d39aca8fd"))
	require.NoError(t, err)

	tests := []struct {
		name          string
		data          *altair.ContributionAndProof
		pubKey        []byte
		domain        [32]byte
		expectedError error
		sig           []byte
	}{
		{
			name: "simple sign",
			data: &altair.ContributionAndProof{
				AggregatorIndex: 7,
				Contribution: &altair.SyncCommitteeContribution{
					Slot:              0,
					BeaconBlockRoot:   [32]byte{},
					SubcommitteeIndex: 0,
					AggregationBits:   []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Signature:         [96]byte{},
				},
				SelectionProof: _byteArray96("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
			},
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: nil,
			sig:           _byteArray("b51e06ac842e499cc34d0d6e5e253e5f48d498d6f075f200da6f65c0a5cd61b37c9b74b3ee7982a94d6a6812cb6a895f0a8ab75380a7df2580663ef5a6d9b477e3b8eecf200ee3875b59859a60e2c6730acd1efb27761af8e8584ceade7617bb"),
		},
		{
			name:          "nil data",
			data:          nil,
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("contrib proof data nil"),
			sig:           _byteArray("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
		},
		{
			name:          "unknown account, should error",
			data:          nil,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
		{
			name:          "nil account, should error",
			data:          nil,
			pubKey:        nil,
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account was not supplied"),
			sig:           nil,
		},
		{
			name:          "empty account, should error",
			data:          nil,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignSyncCommitteeContributionAndProof(test.data, test.domain, test.pubKey)
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

func TestSimpleSigner_SignSyncCommitteeSelectionData(t *testing.T) {
	signer, err := setupNoSlashingProtectionSK(_byteArray("659e875e1b062c03f2f2a57332974d475b97df6cfc581d322e79642d39aca8fd"))
	require.NoError(t, err)

	tests := []struct {
		name          string
		data          *altair.SyncAggregatorSelectionData
		pubKey        []byte
		domain        [32]byte
		expectedError error
		sig           []byte
	}{
		{
			name: "simple sign",
			data: &altair.SyncAggregatorSelectionData{
				Slot:              phase0.Slot(1),
				SubcommitteeIndex: 0,
			},
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: nil,
			sig:           _byteArray("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
		},
		{
			name:          "nil data",
			data:          nil,
			pubKey:        _byteArray("a27c45f7afe6c63363acf886cdad282539fb2cf58b304f2caa95f2ea53048b65a5d41d926c3562e3f18b8b61871375af"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("selection data nil"),
			sig:           _byteArray("a3e966603e64cfd1d091718e3da0e4ed9b13619e7b40d805caf9eadaf84b72dc24fd7f09957a1438f937fbe3e12d6242190dcd5fcbced2b0ef57114ff369c65383eb8561bc56f4ab294ab3a3eba81134e1a90924e85e99e9742009ed4d8f9982"),
		},
		{
			name:          "unknown account, should error",
			data:          nil,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
		{
			name:          "nil account, should error",
			data:          nil,
			pubKey:        nil,
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account was not supplied"),
			sig:           nil,
		},
		{
			name:          "empty account, should error",
			data:          nil,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignSyncCommitteeSelectionData(test.data, test.domain, test.pubKey)
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

func TestSimpleSigner_SignSyncCommittee(t *testing.T) {
	signer, err := setupNoSlashingProtectionSK(_byteArray("15cf88728d14857749044472e5e016c359adf7c020e95e907c438c29fd6f41f9"))
	require.NoError(t, err)

	tests := []struct {
		name          string
		root          []byte
		pubKey        []byte
		domain        [32]byte
		expectedError error
		sig           []byte
	}{
		{
			name:          "simple sign",
			root:          _byteArray("0000000081509579e35e84020ad8751eca180b44df470332d3ad18fc6fd52469"),
			pubKey:        _byteArray("8a90513c2a1ac279aab0c86c9ba6452f809c06d6439a3940aa869fd5cb878c2d7832553faef9059f914b8903c691ef66"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: nil,
			sig:           _byteArray("973523f5a760b8e7d8486970f3eeba3b9a2af13a9b406027cae80967fd8f95a5da605bf073d6901a7ade4fa52be0ae971328fc1010f35d33a4fbca44d2eb7950a3e0010ebef95b1d96091c6952d3b04d65d60765dc37d415222633a2b0afd016"),
		},
		{
			name:          "unknown account, should error",
			root:          _byteArray("0000000081509579e35e84020ad8751eca180b44df470332d3ad18fc6fd52469"),
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
		{
			name:          "nil account, should error",
			root:          _byteArray("0000000081509579e35e84020ad8751eca180b44df470332d3ad18fc6fd52469"),
			pubKey:        nil,
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account was not supplied"),
			sig:           nil,
		},
		{
			name:          "empty account, should error",
			root:          _byteArray("0000000081509579e35e84020ad8751eca180b44df470332d3ad18fc6fd52469"),
			pubKey:        _byteArray(""),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			sig:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := signer.SignSyncCommittee(test.root, test.domain, test.pubKey)
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
