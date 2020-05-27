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
	if options.seed == nil {
		options.GenerateSeed()
	}

	wallet,error := hd.CreateWalletFromSeed(
		options.name,
		options.password,
		options.store,
		options.encryptor,
		options.seed,
	)
	if error != nil {
		return nil, error
	}

	return &KeyVault{
		store:  options.store,
		wallet:	wallet,
	},nil
}