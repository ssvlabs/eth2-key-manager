package dummy

import (
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/encryptor"
	"github.com/google/uuid"
)

type DummyStorage struct{}

func (s *DummyStorage) Name() string                                    { return "" }
func (s *DummyStorage) Network() core.Network                           { return core.MainNetwork }
func (s *DummyStorage) SaveWallet(wallet core.Wallet) error             { return nil }
func (s *DummyStorage) OpenWallet() (core.Wallet, error)                { return nil, nil }
func (s *DummyStorage) ListAccounts() ([]core.ValidatorAccount, error)  { return nil, nil }
func (s *DummyStorage) SaveAccount(account core.ValidatorAccount) error { return nil }
func (s *DummyStorage) OpenAccount(accountId uuid.UUID) (core.ValidatorAccount, error) {
	return nil, nil
}
func (s *DummyStorage) DeleteAccount(accountId uuid.UUID) error                     { return nil }
func (s *DummyStorage) SetEncryptor(encryptor encryptor.Encryptor, password []byte) {}
