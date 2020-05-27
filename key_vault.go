package KeyVault

import (
	hd "github.com/wealdtech/go-eth2-wallet-hd/v2"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type KeyVault struct {
	store wtypes.Store
	wallet wtypes.Wallet
}

func NewKeyVault(options WalletOptions) (*KeyVault,error) {
	options.SetEncryptor(encryptor).SetStore(store).SetWalletName("wallet").SetWalletPassword("password")
	wallet,error := hd.CreateWallet(
		options.name,
		options.password,
		options.store,
		options.encryptor,
	)
	if error != nil {
		return nil, error
	}

	return &KeyVault{
		store:  options.store,
		wallet:	wallet,
	},nil
}