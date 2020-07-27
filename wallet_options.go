package KeyVault

import (
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type WalletOptions struct {
	encryptor wtypes.Encryptor
	password []byte
	storage interface{} // a generic interface as there are a few core storage interfaces (storage, slashing storage and so on)
	seed []byte
}

func (options *WalletOptions) SetEncryptor(encryptor wtypes.Encryptor) *WalletOptions {
	options.encryptor = encryptor
	return options
}

func (options *WalletOptions) SetStorage(storage interface{}) *WalletOptions {
	options.storage = storage
	return options
}

func (options *WalletOptions) SetPassword(password string) *WalletOptions {
	options.password = []byte(password)
	return options
}

func (options *WalletOptions) SetSeed(seed []byte) *WalletOptions {
	options.seed = seed
	return options
}