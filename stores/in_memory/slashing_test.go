package in_memory

import (
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"reflect"
	"testing"
)

func getSlashingStorage() core.SlashingStore {
	return NewInMemStore(
		reflect.TypeOf(KeyVault.KeyVault{}),
		reflect.TypeOf(wallet_hd.HDWallet{}),
		reflect.TypeOf(wallet_hd.HDAccount{}),
	)
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
