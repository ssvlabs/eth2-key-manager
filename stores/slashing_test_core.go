package stores

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"testing"
)

type mockAccount struct {
	id uuid.UUID
	walletid uuid.UUID
}
func (a *mockAccount) ID() uuid.UUID {return a.id}
func (a *mockAccount) WalletID() uuid.UUID {return a.walletid}
func (a *mockAccount) Type() core.AccountType {return core.ValidatorAccount}
func (a *mockAccount) Name() string {return ""}
func (a *mockAccount) PublicKey() e2types.PublicKey {return nil}
func (a *mockAccount) Path() string {return ""}
func (a *mockAccount) Sign(data []byte) (e2types.Signature,error) {return nil,nil}
func (a *mockAccount) SetContext(ctx *core.PortfolioContext){}

func TestingSaveProposal(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name string
		proposal *core.BeaconBlockHeader
		account core.Account
	}{
		{
			name:"simple save",
			proposal: &core.BeaconBlockHeader{
				Slot:          100,
				ProposerIndex: 1,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("A"),
				BodyRoot:      []byte("A"),
			},
			account: &mockAccount{
				id:       uuid.New(),
				walletid: uuid.New(),
			},
		},
	}

	for _,test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveProposal(test.account,test.proposal)
			if err != nil {
				t.Error(err)
				return
			}

			// fetch
			proposal,err := storage.RetrieveProposal(test.account,test.proposal.Slot)
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
		name string
		att *core.BeaconAttestation
		account core.Account
	}{
		{
			name:"simple save",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source:          &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target:          &core.Checkpoint{
					Epoch: 4,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:       uuid.New(),
				walletid: uuid.New(),
			},
		},
		{
			name:"simple save with no change to latest attestation target",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source:          &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target:          &core.Checkpoint{
					Epoch: 3,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:       uuid.New(),
				walletid: uuid.New(),
			},
		},
	}

	for _,test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveAttestation(test.account,test.att)
			if err != nil {
				t.Error(err)
				return
			}

			// fetch
			att,err := storage.RetrieveAttestation(test.account,test.att.Target.Epoch)
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
		id:       uuid.New(),
		walletid: uuid.New(),
	}

	att,err := storage.RetrieveLatestAttestation(account)
	require.NoError(t,err)
	if att != nil {
		t.Errorf("latest attestation should be nil")
		return
	}
}

func TestingSaveLatestAttestation(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name string
		att *core.BeaconAttestation
		account core.Account
	}{
		{
			name:"simple save",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source:          &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target:          &core.Checkpoint{
					Epoch: 4,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:       uuid.New(),
				walletid: uuid.New(),
			},
		},
		{
			name:"simple save with no change to latest attestation target",
			att: &core.BeaconAttestation{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source:          &core.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target:          &core.Checkpoint{
					Epoch: 3,
					Root:  []byte("Root"),
				},
			},
			account: &mockAccount{
				id:       uuid.New(),
				walletid: uuid.New(),
			},
		},
	}

	for _,test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveLatestAttestation(test.account,test.att)
			if err != nil {
				t.Error(err)
				return
			}

			// fetch
			att,err := storage.RetrieveLatestAttestation(test.account)
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
			Source:          &core.Checkpoint{
				Epoch: 1,
				Root:  []byte("Root"),
			},
			Target:          &core.Checkpoint{
				Epoch: 2,
				Root:  []byte("Root"),
			},
		},
		&core.BeaconAttestation{
			Slot:            30,
			CommitteeIndex:  1,
			BeaconBlockRoot: []byte("BeaconBlockRoot"),
			Source:          &core.Checkpoint{
				Epoch: 2,
				Root:  []byte("Root"),
			},
			Target:          &core.Checkpoint{
				Epoch: 3,
				Root:  []byte("Root"),
			},
		},
		&core.BeaconAttestation{
			Slot:            30,
			CommitteeIndex:  1,
			BeaconBlockRoot: []byte("BeaconBlockRoot"),
			Source:          &core.Checkpoint{
				Epoch: 3,
				Root:  []byte("Root"),
			},
			Target:          &core.Checkpoint{
				Epoch: 8,
				Root:  []byte("Root"),
			},
		},
	}
	account := &mockAccount{
		id:       uuid.New(),
		walletid: uuid.New(),
	}

	// save
	for _,att := range attestations {
		err := storage.SaveAttestation(account,att)
		if err != nil {
			t.Error(err)
			return
		}
	}

	tests := []struct{
		name string
		start uint64
		end uint64
		expectedCnt int
	}{
		{
			name: "empty list 1",
			start:0,
			end:1,
			expectedCnt:0,
		},
		{
			name: "empty list 2",
			start:1000,
			end:10010,
			expectedCnt:0,
		},
		{
			name: "simple list 1",
			start:1,
			end:2,
			expectedCnt:1,
		},
		{
			name: "simple list 2",
			start:1,
			end:3,
			expectedCnt:2,
		},
		{
			name: "all",
			start:0,
			end:10,
			expectedCnt:3,
		},
	}

	for _,test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// list
			atts,err := storage.ListAttestations(account,test.start,test.end)
			if err != nil {
				t.Error(err)
				return
			}
			if atts == nil {
				t.Errorf("list attestation returns nil")
				return
			}
			if len(atts) != test.expectedCnt {
				t.Errorf("list attestation returns %d elements, expectd: %d",len(atts), test.expectedCnt)
				return
			}

			// iterate all and compare
			for _, att := range atts {
				if att.Target.Epoch > test.end || att.Source.Epoch < test.start {
					t.Errorf("list attestation returned an element outside what was requested. start: %d end:%d, returned: %d",test.start,test.end,att.Target.Epoch)
					return
				}
			}
		})
	}
}