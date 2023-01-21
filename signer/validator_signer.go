package signer

import (
	"sync"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/google/uuid"
)

type VersionedBeaconBlock struct {
	IsBlinded bool
	Regular   *spec.VersionedBeaconBlock
	Blinded   *api.VersionedBlindedBeaconBlock
}

// ValidatorSigner represents the behavior of the validator signer
type ValidatorSigner interface {
	SignBeaconBlock(block *spec.VersionedBeaconBlock, domain phase0.Domain, pubKey []byte) ([]byte, error)
	//SignBlindedBeaconBlock(block *api.VersionedBlindedBeaconBlock, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignBeaconAttestation(attestation *phase0.AttestationData, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignAggregateAndProof(agg *phase0.AggregateAndProof, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignSlot(slot phase0.Slot, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignEpoch(epoch phase0.Epoch, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignSyncCommittee(msgBlockRoot []byte, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignSyncCommitteeSelectionData(data *altair.SyncAggregatorSelectionData, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignSyncCommitteeContributionAndProof(contribAndProof *altair.ContributionAndProof, domain phase0.Domain, pubKey []byte) ([]byte, error)
	SignRegistration(registration *apiv1.ValidatorRegistration, domain phase0.Domain, pubKey []byte) ([]byte, error)
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
