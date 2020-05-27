package KeyVault

import (
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type WalletOptions struct {
	encryptor wtypes.Encryptor
	password []byte
	walletIndex uint64
	name string
	store wtypes.Store
}

func (options *WalletOptions)SetEncryptor(encryptor wtypes.Encryptor) *WalletOptions {
	options.encryptor = encryptor
	return options
}

func (options *WalletOptions)SetStore(store wtypes.Store) *WalletOptions {
	options.store = store
	return options
}

func (options *WalletOptions)SetWalletName(name string) *WalletOptions {
	options.name = name
	return options
}

func (options *WalletOptions)SetWalletPassword(password string) *WalletOptions {
	options.password = []byte(password)
	return options
}

func (options *WalletOptions)SetWalletIndex(index uint64) *WalletOptions {
	options.walletIndex = index
	return options
}

//func (options *WalletOptions) GenerateSeed() error {
//	seed := make([]byte, 32)
//	_, err := rand.Read(seed)
//
//	options.SetSeed(seed)
//
//	return err
//}