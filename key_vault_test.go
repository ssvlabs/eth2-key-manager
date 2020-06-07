package KeyVault

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/stretchr/testify/assert"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"os"
	"testing"
)

func getStorage() core.PortfolioStorage {
	return in_memory.NewInMemStore()
}

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func portfolio(storage core.PortfolioStorage) (core.Portfolio,error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	options := &PortfolioOptions{}
	options.SetStorage(storage)
	options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"))
	return NewKeyVault(options)
}

func TestMarshalingNoWallets(t *testing.T) {
	p,err := portfolio(getStorage())
	if err != nil {
		t.Error(err)
		return
	}

	m,err := json.Marshal(p)
	if err != nil {
		t.Error(err)
		return
	}

	newP := &KeyVault{}
	err = json.Unmarshal(m,newP)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, p.ID(),newP.ID())
}

func TestMarshalingWithWallets(t *testing.T) {
	p,err := portfolio(getStorage())
	if err != nil {
		t.Error(err)
		return
	}

	w1,err := p.CreateWallet("test1")
	if err != nil {
		t.Error(err)
		return
	}
	w2,err := p.CreateWallet("test2")
	if err != nil {
		t.Error(err)
		return
	}

	m,err := json.Marshal(p)
	if err != nil {
		t.Error(err)
		return
	}

	newP := &KeyVault{}
	err = json.Unmarshal(m,newP)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, p.ID(),newP.ID())
	assert.Equal(t, w1.ID(), newP.indexMapper["test1"])
	assert.Equal(t, w2.ID(), newP.indexMapper["test2"])
}
