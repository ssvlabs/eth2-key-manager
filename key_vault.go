package KeyVault

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/slashing_protectors"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"github.com/wealdtech/go-indexer"
)

// This is an EIP 2333,2334,2335 compliant hierarchical deterministic portfolio
//https://eips.ethereum.org/EIPS/eip-2333
//https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
//https://eips.ethereum.org/EIPS/eip-2335
type KeyVault struct {
	Storage  interface{}
	enableSimpleSigner bool
	walletsIndexer indexer.Index // maps indexs <> names
	walletIds []uuid.UUID
	context *core.PortfolioContext
	key *core.EncryptableSeed
}

//func OpenKeyVault(options *WalletOptions) (*KeyVault,error) {
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
//		slashingProtection := slashing_protectors.NewNormalProtection(options.store.(slashing_protectors.SlashingStore))
//		signer = validator_signer.NewSimpleSigner(wallet,slashingProtection)
//	}
//
//	return &KeyVault{
//		Store:  options.store,
//		Wallet: wallet,
//		Signer: signer,
//	},nil
//}

func NewKeyVault(options *WalletOptions) (*KeyVault,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,err
	}

	// set encryptor
	if options.encryptor == nil {
		options.setNoEncryptor()
	}

	// set seed
	var seed *core.EncryptableSeed
	if options.seed == nil {
		options.GenerateSeed()
	}
	seed = core.NewEncryptableSeed(options.seed, options.encryptor)

	// storage
	if _,ok := options.storage.(core.PortfolioStorage); !ok {
		return nil,fmt.Errorf("storage does not implement PortfolioStorage")
	}

	// signer
	if options.enableSimpleSigner {
		if _,ok := options.storage.(slashing_protectors.SlashingStore); !ok {
			return nil,fmt.Errorf("storage does not implement SlashingStore")
		}
	}

	// lock policy
	if options.portfolioLockPolicy == nil {
		options.setNoLockPolicy()
	}

	// portfolio context
	context := &core.PortfolioContext {
		Storage: 		options.storage.(core.PortfolioStorage),
		Encryptor: 		options.encryptor,
		LockPolicy:		options.portfolioLockPolicy,
		LockPassword:	options.password,
	}

	return &KeyVault{
		Storage:            options.storage,
		enableSimpleSigner: options.enableSimpleSigner,
		walletsIndexer:     indexer.Index{},
		walletIds:          make([]uuid.UUID,0),
		context:            context,
		key:				seed,
	},nil

	//var wallet wtypes.Wallet
	//var error error
	//if options.seed != nil {
	//	wallet,error = hd.CreateWalletFromSeed(
	//		options.name,
	//		options.password,
	//		options.store,
	//		options.encryptor,
	//		options.seed,
	//	)
	//	if error != nil {
	//		return nil, error
	//	}
	//} else {
	//	wallet,error = hd.CreateWallet(
	//		options.name,
	//		options.password,
	//		options.store,
	//		options.encryptor,
	//	)
	//	if error != nil {
	//		return nil, error
	//	}
	//}
	//
	//var signer validator_signer.ValidatorSigner
	//if options.enableSimpleSigner{
	//	slashingProtection := slashing_protectors.NewNormalProtection(options.store.(slashing_protectors.SlashingStore))
	//	signer = validator_signer.NewSimpleSigner(wallet,slashingProtection)
	//}
	//
	//return &KeyVault{
	//	Store:  options.store,
	//	Wallet: wallet,
	//	Signer: signer,
	//},nil
}
