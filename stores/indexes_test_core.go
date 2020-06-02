package stores

import (
	"bytes"
	"github.com/google/uuid"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
	"testing"
)

func TestingAccountIndexes(storage wtypes.Store, t *testing.T) {
	walletId := uuid.New()
	err := storage.StoreAccountsIndex(walletId,[]byte("test"))
	if err != nil {
		t.Error(err)
		return
	}

	data,err := storage.RetrieveAccountsIndex(walletId)
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(data,[]byte("test")) != 0 {
		t.Errorf("indexes were not fetched correctly")
	}
}