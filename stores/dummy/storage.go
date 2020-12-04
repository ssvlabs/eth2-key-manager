package dummy

import (
	"github.com/google/uuid"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/encryptor"
)

type Storage struct{}

// Name does nothing
func (s *Storage) Name() string { return "" }

// Network does nothing
func (s *Storage) Network() core.Network { return core.MainNetwork }

// SaveWallet does nothing
func (s *Storage) SaveWallet(_ core.Wallet) error { return nil }

// OpenWallet does nothing
func (s *Storage) OpenWallet() (core.Wallet, error) { return nil, nil }

// ListAccounts does nothing
func (s *Storage) ListAccounts() ([]core.ValidatorAccount, error) { return nil, nil }

// SaveAccount nothing
func (s *Storage) SaveAccount(_ core.ValidatorAccount) error { return nil }

// OpenAccount does nothing
func (s *Storage) OpenAccount(_ uuid.UUID) (core.ValidatorAccount, error) {
	return nil, nil
}

// DeleteAccount does nothing
func (s *Storage) DeleteAccount(_ uuid.UUID) error { return nil }

// SetEncryptor does nothing
func (s *Storage) SetEncryptor(_ encryptor.Encryptor, _ []byte) {}
