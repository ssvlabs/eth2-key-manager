package in_memory

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func getSlashingStorage() core.SlashingStore {
	return NewInMemStore()
}

func TestSavingProposal (t *testing.T) {
	stores.TestingSaveProposal(getSlashingStorage(),t)
}

func TestSavingAttestation (t *testing.T) {
	stores.TestingSaveAttestation(getSlashingStorage(),t)
}

func TestSavingLatestAttestation (t *testing.T) {
	stores.TestingSaveLatestAttestation(getSlashingStorage(),t)
}

func TestListingAttestation (t *testing.T) {
	stores.TestingListingAttestation(getSlashingStorage(),t)
}
