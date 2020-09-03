package in_memory

import (
	"testing"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores"
)

func getSlashingStorage() core.SlashingStore {
	return NewInMemStore()
}

func TestSavingProposal(t *testing.T) {
	stores.TestingSaveProposal(getSlashingStorage(), t)
}

func TestSavingAttestation(t *testing.T) {
	stores.TestingSaveAttestation(getSlashingStorage(), t)
}

func TestSavingLatestAttestation(t *testing.T) {
	stores.TestingSaveLatestAttestation(getSlashingStorage(), t)
}

func TestRetrieveEmptyLatestAttestation(t *testing.T) {
	stores.TestingRetrieveEmptyLatestAttestation(getSlashingStorage(), t)
}

func TestListingAttestation(t *testing.T) {
	stores.TestingListingAttestation(getSlashingStorage(), t)
}
