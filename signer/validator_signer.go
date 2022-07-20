package signer

import (
	"sync"

	"github.com/prysmaticlabs/prysm/consensus-types/interfaces"

	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"

	"github.com/google/uuid"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"

	"github.com/bloxapp/eth2-key-manager/core"
)

// ValidatorSigner represents the behavior of the validator signer
type ValidatorSigner interface {
	SignBeaconBlock(block interfaces.BeaconBlock, domain []byte, pubKey []byte) ([]byte, error)
	SignBeaconAttestation(attestation *eth.AttestationData, domain []byte, pubKey []byte) ([]byte, error)
	SignAggregateAndProof(agg *eth.AggregateAttestationAndProof, domain []byte, pubKey []byte) ([]byte, error)
	SignSlot(slot types.Slot, domain []byte, pubKey []byte) ([]byte, error)
	SignEpoch(epoch types.Epoch, domain []byte, pubKey []byte) ([]byte, error)
	SignSyncCommittee(msgBlockRoot []byte, domain []byte, pubKey []byte) ([]byte, error)
	SignSyncCommitteeSelectionData(data *eth.SyncAggregatorSelectionData, domain []byte, pubKey []byte) ([]byte, error)
	SignSyncCommitteeContributionAndProof(contribAndProof *eth.ContributionAndProof, domain []byte, pubKey []byte) ([]byte, error)
	SignRegistration(registration *eth.ValidatorRegistrationV1, domain []byte, pubKey []byte) ([]byte, error)
}

// SimpleSigner implements ValidatorSigner interface
type SimpleSigner struct {
	wallet            core.Wallet
	slashingProtector core.SlashingProtector
	network           core.Network
	signLocks         map[string]*sync.RWMutex
	mapLock           *sync.RWMutex
}

// NewSimpleSigner is the constructor of SimpleSigner
func NewSimpleSigner(wallet core.Wallet, slashingProtector core.SlashingProtector, network core.Network) *SimpleSigner {
	return &SimpleSigner{
		wallet:            wallet,
		slashingProtector: slashingProtector,
		network:           network,
		signLocks:         map[string]*sync.RWMutex{},
		mapLock:           &sync.RWMutex{},
	}
}

// lock locks signer
func (signer *SimpleSigner) lock(accountID uuid.UUID, operation string) {
	signer.mapLock.Lock()
	defer signer.mapLock.Unlock()

	k := accountID.String() + "_" + operation
	if val, ok := signer.signLocks[k]; ok {
		val.Lock()
	} else {
		signer.signLocks[k] = &sync.RWMutex{}
		signer.signLocks[k].Lock()
	}
}

func (signer *SimpleSigner) unlock(accountID uuid.UUID, operation string) {
	signer.mapLock.RLock()
	defer signer.mapLock.RUnlock()

	k := accountID.String() + "_" + operation
	if val, ok := signer.signLocks[k]; ok {
		val.Unlock()
	}
}
