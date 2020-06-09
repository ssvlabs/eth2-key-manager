package KeyVault

import (
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type PortfolioOptions struct {
	encryptor wtypes.Encryptor
	password []byte
	storage interface{} // a generic interface as there are a few core storage interfaces (storage, slashing storage and so on)
	seed []byte
}

func (options *PortfolioOptions)SetEncryptor(encryptor wtypes.Encryptor) *PortfolioOptions {
	options.encryptor = encryptor
	return options
}

func (options *PortfolioOptions)SetStorage(storage interface{}) *PortfolioOptions {
	options.storage = storage
	return options
}

func (options *PortfolioOptions)SetPassword(password string) *PortfolioOptions {
	options.password = []byte(password)
	return options
}

func (options *PortfolioOptions)SetSeed(seed []byte) *PortfolioOptions {
	options.seed = seed
	return options
}