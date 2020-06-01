package hashicorp

import (
	"context"
	slash "github.com/bloxapp/KeyVault/slashing_protectors"
	"github.com/bloxapp/KeyVault/stores"
	"github.com/hashicorp/vault/sdk/logical"
	"testing"
)

func getSlashingStorage() slash.SlashingStore {
	return NewHashicorpVaultStore(&logical.InmemStorage{},context.Background())
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
