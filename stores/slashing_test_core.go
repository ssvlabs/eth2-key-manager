package stores

import (
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
)

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

type mockAccount struct {
	id            uuid.UUID
	validationKey *big.Int
}

func (a *mockAccount) ID() uuid.UUID    { return a.id }
func (a *mockAccount) Name() string     { return "" }
func (a *mockAccount) BasePath() string { return "" }
func (a *mockAccount) ValidatorPublicKey() e2types.PublicKey {
	priv, _ := e2types.BLSPrivateKeyFromBytes(a.validationKey.Bytes())
	return priv.PublicKey()
}
func (a *mockAccount) WithdrawalPublicKey() e2types.PublicKey                   { return nil }
func (a *mockAccount) ValidationKeySign(data []byte) (e2types.Signature, error) { return nil, nil }
func (a *mockAccount) WithdrawalKeySign(data []byte) (e2types.Signature, error) { return nil, nil }
func (a *mockAccount) GetDepositData() (map[string]interface{}, error)          { return nil, nil }
func (a *mockAccount) SetContext(ctx *core.WalletContext)                       {}

func TestingSaveProposal(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name     string
		proposal *core.BeaconBlockHeader
		account  core.ValidatorAccount
	}{
		{
			name: "simple save",
			proposal: &core.BeaconBlockHeader{
				Slot:          100,
				ProposerIndex: 1,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("A"),
				BodyRoot:      []byte("A"),
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveProposal(test.account.ValidatorPublicKey(), test.proposal)
			require.NoError(t, err)

			// fetch
			proposal, err := storage.RetrieveProposal(test.account.ValidatorPublicKey(), test.proposal.Slot)
			require.NoError(t, err)
			require.NotNil(t, proposal)
			require.True(t, proposal.Compare(test.proposal))
		})
	}
}

func TestingSaveAttestation(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name    string
		att     *core.BeaconAttestation
		account core.ValidatorAccount
	}{
		{
			name: "simple save",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &core.Checkpoint{
					Epoch: 4,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
		{
			name: "simple save with no change to latest attestation target",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &core.Checkpoint{
					Epoch: 3,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveHighestAttestation(test.account.ValidatorPublicKey(), test.att)
			require.NoError(t, err)

			// fetch
			att := storage.RetrieveHighestAttestation(test.account.ValidatorPublicKey())
			require.NotNil(t, att)
			require.True(t, att.Compare(test.att))
		})
	}
}

//func TestingRetrieveEmptyLatestAttestation(storage core.SlashingStore, t *testing.T) {
//	account := &mockAccount{
//		id:            uuid.New(),
//		validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
//	}
//
//	att, err := storage.RetrieveLatestAttestation(account.ValidatorPublicKey())
//	require.NoError(t, err)
//	require.Nil(t, att)
//}

func TestingSaveHighestAttestation(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name    string
		att     *core.BeaconAttestation
		account core.ValidatorAccount
	}{
		{
			name: "simple save",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &core.Checkpoint{
					Epoch: 4,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
		{
			name: "simple save with no change to latest attestation target",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &core.Checkpoint{
					Epoch: 3,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveHighestAttestation(test.account.ValidatorPublicKey(), test.att)
			require.NoError(t, err)

			// fetch
			att := storage.RetrieveHighestAttestation(test.account.ValidatorPublicKey())
			require.NotNil(t, att)
			require.True(t, att.Compare(test.att))
		})
	}
}

//func TestingListingAttestation(storage core.SlashingStore, t *testing.T) {
//	attestations := []*core.BeaconAttestation{
//		&core.BeaconAttestation{
//			Slot:            30,
//			CommitteeIndex:  1,
//			BeaconBlockRoot: []byte("BeaconBlockRoot"),
//			Source: &core.Checkpoint{
//				Epoch: 1,
//				Root:  []byte("Root"),
//			},
//			Target: &core.Checkpoint{
//				Epoch: 2,
//				Root:  []byte("Root"),
//			},
//		},
//		&core.BeaconAttestation{
//			Slot:            30,
//			CommitteeIndex:  1,
//			BeaconBlockRoot: []byte("BeaconBlockRoot"),
//			Source: &core.Checkpoint{
//				Epoch: 2,
//				Root:  []byte("Root"),
//			},
//			Target: &core.Checkpoint{
//				Epoch: 3,
//				Root:  []byte("Root"),
//			},
//		},
//		&core.BeaconAttestation{
//			Slot:            30,
//			CommitteeIndex:  1,
//			BeaconBlockRoot: []byte("BeaconBlockRoot"),
//			Source: &core.Checkpoint{
//				Epoch: 3,
//				Root:  []byte("Root"),
//			},
//			Target: &core.Checkpoint{
//				Epoch: 8,
//				Root:  []byte("Root"),
//			},
//		},
//	}
//	account := &mockAccount{
//		id:            uuid.New(),
//		validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
//	}
//
//	// save
//	for _, att := range attestations {
//		err := storage.SaveHighestAttestation(account.ValidatorPublicKey(), att)
//		require.NoError(t, err)
//	}
//
//	tests := []struct {
//		name        string
//		start       uint64
//		end         uint64
//		expectedCnt int
//	}{
//		{
//			name:        "empty list 1",
//			start:       0,
//			end:         1,
//			expectedCnt: 0,
//		},
//		{
//			name:        "empty list 2",
//			start:       1000,
//			end:         10010,
//			expectedCnt: 0,
//		},
//		{
//			name:        "simple list 1",
//			start:       1,
//			end:         2,
//			expectedCnt: 1,
//		},
//		{
//			name:        "simple list 2",
//			start:       1,
//			end:         3,
//			expectedCnt: 2,
//		},
//		{
//			name:        "all",
//			start:       0,
//			end:         10,
//			expectedCnt: 3,
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			// list
//			atts, err := storage.ListAttestations(account.ValidatorPublicKey(), test.start, test.end)
//			require.NoError(t, err)
//			require.NotNil(t, atts)
//			require.Len(t, atts, test.expectedCnt)
//
//			// iterate all and compare
//			for _, att := range atts {
//				require.False(t, att.Target.Epoch > test.end || att.Source.Epoch < test.start)
//			}
//		})
//	}
//}
