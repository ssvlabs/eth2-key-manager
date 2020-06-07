package KeyVault

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/slashing_protectors"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// This is an EIP 2333,2334,2335 compliant hierarchical deterministic portfolio
//https://eips.ethereum.org/EIPS/eip-2333
//https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
//https://eips.ethereum.org/EIPS/eip-2335
type KeyVault struct {
	id uuid.UUID
	Storage  interface{}
	enableSimpleSigner bool
	indexMapper map[string]uuid.UUID
	context *core.PortfolioContext
	key *core.DerivableKey
}


func (vault *KeyVault) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = vault.id
	data["enableSimpleSigner"] = vault.enableSimpleSigner
	data["indexMapper"] = vault.indexMapper

	return json.Marshal(data)
}

func (vault *KeyVault) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	var err error

	// id
	if val, exists := v["id"]; exists {
		vault.id,err = uuid.Parse(val.(string))
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: id")}

	// simple signer
	if val, exists := v["enableSimpleSigner"]; exists {
		vault.enableSimpleSigner = val.(bool)
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: enableSimpleSigner")}

	// indexMapper
	if val, exists := v["indexMapper"]; exists {
		vault.indexMapper = make(map[string]uuid.UUID)
		for k,v := range val.(map[string]interface{}) {
			vault.indexMapper[k],err = uuid.Parse(v.(string))
			if err != nil {
				return err
			}
		}
	} else {return fmt.Errorf("could not find var: indexMapper")}
	return nil
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

func NewKeyVault(options *PortfolioOptions) (*KeyVault,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,err
	}

	// set encryptor
	if options.encryptor == nil {
		options.setNoEncryptor()
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
	}

	// signer
	if options.enableSimpleSigner {
		if _,ok := options.storage.(slashing_protectors.SlashingStore); !ok {
			return nil,fmt.Errorf("storage does not implement SlashingStore")
		}
	}

	// portfolio context
	context := &core.PortfolioContext {
		Storage: 		options.storage.(core.PortfolioStorage),
	}

	ret := &KeyVault{
		Storage:            options.storage,
		enableSimpleSigner: options.enableSimpleSigner,
		indexMapper:        make(map[string]uuid.UUID),
		context:            context,
		key:				seed,
	}

	// update context with portfolio id
	context.PortfolioId = ret.ID()

	return ret,nil
}
