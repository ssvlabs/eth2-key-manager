package KeyVault

import (
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type KeyVaultOptions struct {
	encryptor wtypes.Encryptor
	password []byte
	storage interface{} // a generic interface as there are a few core storage interfaces (storage, slashing storage and so on)
	seed []byte
}

func (options *KeyVaultOptions) SetEncryptor(encryptor wtypes.Encryptor) *KeyVaultOptions {
	options.encryptor = encryptor
	return options
}

func (options *KeyVaultOptions) SetStorage(storage interface{}) *KeyVaultOptions {
	options.storage = storage
	return options
}

func (options *KeyVaultOptions) SetPassword(password string) *KeyVaultOptions {
	options.password = []byte(password)
	return options
}

func (options *KeyVaultOptions) SetSeed(seed []byte) *KeyVaultOptions {
	options.seed = seed
	return options
}