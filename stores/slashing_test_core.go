package stores

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"math/big"
	"testing"
)

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

type mockAccount struct {
	id            uuid.UUID
	validationKey *big.Int
}

func (a *mockAccount) ID() uuid.UUID { return a.id }
func (a *mockAccount) Name() string  { return "" }
func (a *mockAccount) ValidatorPublicKey() e2types.PublicKey {
	priv, _ := e2types.BLSPrivateKeyFromBytes(a.validationKey.Bytes())
	return priv.PublicKey()
}
func (a *mockAccount) WithdrawalPublicKey() e2types.PublicKey                   { return nil }
func (a *mockAccount) ValidationKeySign(data []byte) (e2types.Signature, error) { return nil, nil }
func (a *mockAccount) WithdrawalKeySign(data []byte) (e2types.Signature, error) { return nil, nil }
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
			if err != nil {
				t.Error(err)
				return
			}

			// fetch
			proposal, err := storage.RetrieveProposal(test.account.ValidatorPublicKey(), test.proposal.Slot)
			if err != nil {
				t.Error(err)
				return
			}
			if proposal == nil {
				t.Errorf("proposal not saved and retrieved")
				return
			}
			if proposal.Compare(test.proposal) != true {
				t.Errorf("retrieved proposal not matching saved attestation")
				return
			}
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
			err := storage.SaveAttestation(test.account.ValidatorPublicKey(), test.att)
			if err != nil {
				t.Error(err)
				return
			}

			// fetch
			att, err := storage.RetrieveAttestation(test.account.ValidatorPublicKey(), test.att.Target.Epoch)
			if err != nil {
				t.Error(err)
				return
			}
			if att == nil {
				t.Errorf("attestation not saved and retrieved")
				return
			}
			if att.Compare(test.att) != true {
				t.Errorf("retrieved attestation not matching saved attestation")
				return
			}
		})
	}
}

func TestingRetrieveEmptyLatestAttestation(storage core.SlashingStore, t *testing.T) {
	account := &mockAccount{
		id:            uuid.New(),
		validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
	}

	att, err := storage.RetrieveLatestAttestation(account.ValidatorPublicKey())
	require.NoError(t, err)
	if att != nil {
		t.Errorf("latest attestation should be nil")
		return
	}
}

func TestingSaveLatestAttestation(storage core.SlashingStore, t *testing.T) {
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
			err := storage.SaveLatestAttestation(test.account.ValidatorPublicKey(), test.att)
			if err != nil {
				t.Error(err)
				return
			}

			// fetch
			att, err := storage.RetrieveLatestAttestation(test.account.ValidatorPublicKey())
			if err != nil {
				t.Error(err)
				return
			}
			if att == nil {
				t.Errorf("latest attestation not saved and retrieved")
				return
			}
			if att.Compare(test.att) != true {
				t.Errorf("retrieved latest attestation not matching saved attestation")
				return
			}
		})
	}
}

func TestingListingAttestation(storage core.SlashingStore, t *testing.T) {
	attestations := []*core.BeaconAttestation{
		&core.BeaconAttestation{
			Slot:            30,
			CommitteeIndex:  1,
			BeaconBlockRoot: []byte("BeaconBlockRoot"),
			Source: &core.Checkpoint{
				Epoch: 1,
				Root:  []byte("Root"),
			},
			Target: &core.Checkpoint{
				Epoch: 2,
				Root:  []byte("Root"),
			},
		},
		&core.BeaconAttestation{
			Slot:            30,
			CommitteeIndex:  1,
			BeaconBlockRoot: []byte("BeaconBlockRoot"),
			Source: &core.Checkpoint{
				Epoch: 2,
				Root:  []byte("Root"),
			},
			Target: &core.Checkpoint{
				Epoch: 3,
				Root:  []byte("Root"),
			},
		},
		&core.BeaconAttestation{
			Slot:            30,
			CommitteeIndex:  1,
			BeaconBlockRoot: []byte("BeaconBlockRoot"),
			Source: &core.Checkpoint{
				Epoch: 3,
				Root:  []byte("Root"),
			},
			Target: &core.Checkpoint{
				Epoch: 8,
				Root:  []byte("Root"),
			},
		},
	}
	account := &mockAccount{
		id:            uuid.New(),
		validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
	}

	// save
	for _, att := range attestations {
		err := storage.SaveAttestation(account.ValidatorPublicKey(), att)
		if err != nil {
			t.Error(err)
			return
		}
	}

	tests := []struct {
		name        string
		start       uint64
		end         uint64
		expectedCnt int
	}{
		{
			name:        "empty list 1",
			start:       0,
			end:         1,
			expectedCnt: 0,
		},
		{
			name:        "empty list 2",
			start:       1000,
			end:         10010,
			expectedCnt: 0,
		},
		{
			name:        "simple list 1",
			start:       1,
			end:         2,
			expectedCnt: 1,
		},
		{
			name:        "simple list 2",
			start:       1,
			end:         3,
			expectedCnt: 2,
		},
		{
			name:        "all",
			start:       0,
			end:         10,
			expectedCnt: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// list
			atts, err := storage.ListAttestations(account.ValidatorPublicKey(), test.start, test.end)
			if err != nil {
				t.Error(err)
				return
			}
			if atts == nil {
				t.Errorf("list attestation returns nil")
				return
			}
			if len(atts) != test.expectedCnt {
				t.Errorf("list attestation returns %d elements, expectd: %d", len(atts), test.expectedCnt)
				return
			}

			// iterate all and compare
			for _, att := range atts {
				if att.Target.Epoch > test.end || att.Source.Epoch < test.start {
					t.Errorf("list attestation returned an element outside what was requested. start: %d end:%d, returned: %d", test.start, test.end, att.Target.Epoch)
					return
				}
			}
		})
	}
}
