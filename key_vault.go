package KeyVault

import (
	"github.com/bloxapp/KeyVault/ValidatorSigner"
	"github.com/bloxapp/KeyVault/slashing_protectors"
	hd "github.com/wealdtech/go-eth2-wallet-hd/v2"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type KeyVault struct {
	store wtypes.Store
	wallet wtypes.Wallet
	signer ValidatorSigner.ValidatorSigner
}

func NewKeyVault(options WalletOptions) (*KeyVault,error) {
	wallet,error := hd.CreateWallet(
		options.name,
		options.password,
		options.store,
		options.encryptor,
	)
	if error != nil {
		return nil, error
	}

	var signer ValidatorSigner.ValidatorSigner
	if options.enableSimpleSigner{
		slashingProtection := slashing_protectors.NewNormalProtection(options.store)
		signer = ValidatorSigner.NewSimpleSigner(wallet,slashingProtection)
	}

	return &KeyVault{
		store:  options.store,
		wallet:	wallet,
		signer: signer,
	},nil
}