package in_memory

import (
	slash "github.com/bloxapp/KeyVault/slashing_protectors"
	"github.com/bloxapp/KeyVault/stores"
	"testing"
)

func getSlashingStorage() slash.SlashingStore {
	return NewInMemStore()
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
