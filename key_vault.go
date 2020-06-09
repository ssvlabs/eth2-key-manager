package KeyVault

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// This is an EIP 2333,2334,2335 compliant hierarchical deterministic portfolio
//https://eips.ethereum.org/EIPS/eip-2333
//https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
//https://eips.ethereum.org/EIPS/eip-2335
type KeyVault struct {
	id                 uuid.UUID
	enableSimpleSigner bool
	indexMapper        map[string]uuid.UUID
	Context            *core.PortfolioContext
	key                *core.DerivableKey
}

//func OpenKeyVault(options *PortfolioOptions) (*KeyVault,error) {
//	if err := e2types.InitBLS(); err != nil { // very important!
//		return nil,err
//	}
//
//	wallet,err := hd.OpenWallet(
//			options.name,
//			options.store,
//			options.encryptor,
//		)
//	if err != nil {
//		return nil, err
//	}
//
//	var signer validator_signer.ValidatorSigner
//	if options.enableSimpleSigner{
//		slashingProtection := slashing_protection.NewNormalProtection(options.store.(slashing_protection.SlashingStore))
//		signer = validator_signer.NewSimpleSigner(wallet,slashingProtection)
//	}
//
//	return &KeyVault{
//		Store:  options.store,
//		Wallet: wallet,
//		Signer: signer,
//	},nil
//}

func NewKeyVault(options *PortfolioOptions) (*KeyVault,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,err
	}

	// set seed
	if options.seed == nil {
		options.GenerateSeed()
	}
	seed,err := core.BaseKeyFromSeed(options.seed)
	if err != nil {
		return nil,err
	}

	// storage
	if _,ok := options.storage.(core.PortfolioStorage); !ok {
		return nil,fmt.Errorf("storage does not implement PortfolioStorage")
	} else {
		if options.encryptor != nil && options.password != nil {
			options.storage.(core.PortfolioStorage).SetEncryptor(options.encryptor,options.password)
		}
	}

	// signer
	if options.enableSimpleSigner {
		if _,ok := options.storage.(core.SlashingStore); !ok {
			return nil,fmt.Errorf("storage does not implement SlashingStore")
		}
	}

	// portfolio Context
	context := &core.PortfolioContext {
		Storage:	options.storage.(core.PortfolioStorage),
	}

	ret := &KeyVault{
		enableSimpleSigner: options.enableSimpleSigner,
		indexMapper:        make(map[string]uuid.UUID),
		Context:            context,
		key:                seed,
	}

	// update Context with portfolio id
	context.PortfolioId = ret.ID()

	return ret,nil
}
