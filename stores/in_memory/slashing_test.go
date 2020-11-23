package in_memory

import (
	"testing"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores"
)

func getSlashingStorage() core.SlashingStore {
	return NewInMemStore(core.MainNetwork)
}

func TestSavingProposal(t *testing.T) {
	stores.TestingSaveProposal(getSlashingStorage(), t)
}

func TestSavingAttestation(t *testing.T) {
	stores.TestingSaveAttestation(getSlashingStorage(), t)
}

func TestSavingHighestAttestation(t *testing.T) {
	stores.TestingSaveHighestAttestation(getSlashingStorage(), t)
}

//func TestRetrieveEmptyLatestAttestation(t *testing.T) {
//	stores.TestingRetrieveEmptyLatestAttestation(getSlashingStorage(), t)
//}
//
//func TestListingAttestation(t *testing.T) {
//	stores.TestingListingAttestation(getSlashingStorage(), t)
//}
