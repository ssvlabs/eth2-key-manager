package dummy

import (
	"github.com/google/uuid"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/encryptor"
)

type Storage struct{}

func (s *Storage) Name() string                                   { return "" }
func (s *Storage) Network() core.Network                          { return core.MainNetwork }
func (s *Storage) SaveWallet(_ core.Wallet) error                 { return nil }
func (s *Storage) OpenWallet() (core.Wallet, error)               { return nil, nil }
func (s *Storage) ListAccounts() ([]core.ValidatorAccount, error) { return nil, nil }
func (s *Storage) SaveAccount(_ core.ValidatorAccount) error      { return nil }
func (s *Storage) OpenAccount(_ uuid.UUID) (core.ValidatorAccount, error) {
	return nil, nil
}
func (s *Storage) DeleteAccount(_ uuid.UUID) error              { return nil }
func (s *Storage) SetEncryptor(_ encryptor.Encryptor, _ []byte) {}
