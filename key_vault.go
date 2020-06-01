package KeyVault

import (
	"github.com/bloxapp/KeyVault/slashing_protectors"
	"github.com/bloxapp/KeyVault/validator_signer"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	hd "github.com/wealdtech/go-eth2-wallet-hd/v2"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type KeyVault struct {
	Store  wtypes.Store
	Wallet wtypes.Wallet
	Signer validator_signer.ValidatorSigner
}

func OpenKeyVault(options *WalletOptions) (*KeyVault,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,err
	}

	wallet,err := hd.OpenWallet(
			options.name,
			options.store,
			options.encryptor,
		)
	if err != nil {
		return nil, err
	}

	var signer validator_signer.ValidatorSigner
	if options.enableSimpleSigner{
		slashingProtection := slashing_protectors.NewNormalProtection(options.store.(slashing_protectors.SlashingStore))
		signer = validator_signer.NewSimpleSigner(wallet,slashingProtection)
	}

	return &KeyVault{
		Store:  options.store,
		Wallet: wallet,
		Signer: signer,
	},nil
}

func NewKeyVault(options *WalletOptions) (*KeyVault,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,err
	}

	var wallet wtypes.Wallet
	var error error
	if options.seed != nil {
		wallet,error = hd.CreateWalletFromSeed(
			options.name,
			options.password,
			options.store,
			options.encryptor,
			options.seed,
		)
		if error != nil {
			return nil, error
		}
	} else {
		wallet,error = hd.CreateWallet(
			options.name,
			options.password,
			options.store,
			options.encryptor,
		)
		if error != nil {
			return nil, error
		}
	}

	var signer validator_signer.ValidatorSigner
	if options.enableSimpleSigner{
		slashingProtection := slashing_protectors.NewNormalProtection(options.store.(slashing_protectors.SlashingStore))
		signer = validator_signer.NewSimpleSigner(wallet,slashingProtection)
	}

	return &KeyVault{
		Store:  options.store,
		Wallet: wallet,
		Signer: signer,
	},nil
}