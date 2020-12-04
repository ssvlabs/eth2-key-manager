package eth2keymanager

import (
	"github.com/bloxapp/eth2-key-manager/core"
	encryptor2 "github.com/bloxapp/eth2-key-manager/encryptor"
)

// KeyVaultOptions contains options to create a new key vault object
type KeyVaultOptions struct {
	encryptor  encryptor2.Encryptor
	password   []byte
	storage    interface{} // a generic interface as there are a few core storage interfaces (storage, slashing storage and so on)
	walletType core.WalletType
}

// SetEncryptor is the encryptor setter
func (options *KeyVaultOptions) SetEncryptor(encryptor encryptor2.Encryptor) *KeyVaultOptions {
	options.encryptor = encryptor
	return options
}

// SetStorage is the storage setter
func (options *KeyVaultOptions) SetStorage(storage interface{}) *KeyVaultOptions {
	options.storage = storage
	return options
}

// SetPassword is the password setter
func (options *KeyVaultOptions) SetPassword(password string) *KeyVaultOptions {
	options.password = []byte(password)
	return options
}

// SetWalletType is the wallet type setter
func (options *KeyVaultOptions) SetWalletType(walletType core.WalletType) *KeyVaultOptions {
	options.walletType = walletType
	return options
}
