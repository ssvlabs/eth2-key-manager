package stores

import (
	"math/big"
	"testing"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

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
func (a *mockAccount) ValidatorPublicKey() []byte {
	priv, _ := e2types.BLSPrivateKeyFromBytes(a.validationKey.Bytes())
	return priv.PublicKey().Marshal()
}
func (a *mockAccount) WithdrawalPublicKey() []byte                     { return nil }
func (a *mockAccount) ValidationKeySign(data []byte) ([]byte, error)   { return nil, nil }
func (a *mockAccount) GetDepositData() (map[string]interface{}, error) { return nil, nil }
func (a *mockAccount) SetContext(ctx *core.WalletContext)              {}

func TestingSaveProposal(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name     string
		proposal *eth.BeaconBlock
		account  core.ValidatorAccount
	}{
		{
			name: "simple save",
			proposal: &eth.BeaconBlock{
				Slot:          100,
				ProposerIndex: 1,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("A"),
				Body:          &eth.BeaconBlockBody{},
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

			// test equal
			aRoot, err := proposal.HashTreeRoot()
			require.NoError(t, err)
			bRoot, err := proposal.HashTreeRoot()
			require.NoError(t, err)
			require.EqualValues(t, aRoot, bRoot)
		})
	}
}

func TestingSaveAttestation(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name    string
		att     *eth.AttestationData
		account core.ValidatorAccount
	}{
		{
			name: "simple save",
			att: &eth.AttestationData{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &eth.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &eth.Checkpoint{
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
			att: &eth.AttestationData{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &eth.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &eth.Checkpoint{
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

			// test equal
			aRoot, err := att.HashTreeRoot()
			require.NoError(t, err)
			bRoot, err := test.att.HashTreeRoot()
			require.NoError(t, err)
			require.EqualValues(t, aRoot, bRoot)
		})
	}
}

func TestingSaveHighestAttestation(storage core.SlashingStore, t *testing.T) {
	tests := []struct {
		name    string
		att     *eth.AttestationData
		account core.ValidatorAccount
	}{
		{
			name: "simple save",
			att: &eth.AttestationData{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &eth.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &eth.Checkpoint{
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
			att: &eth.AttestationData{
				Slot:            30,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &eth.Checkpoint{
					Epoch: 1,
					Root:  []byte("Root"),
				},
				Target: &eth.Checkpoint{
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

			// test equal
			aRoot, err := att.HashTreeRoot()
			require.NoError(t, err)
			bRoot, err := test.att.HashTreeRoot()
			require.NoError(t, err)
			require.EqualValues(t, aRoot, bRoot)
		})
	}
}
