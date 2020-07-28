package wallet_hd

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
)

// according to https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
const (
	BaseAccountPath = "/%d"
	WithdrawalKeyPath = BaseAccountPath + "/0"
	ValidatorKeyPath = WithdrawalKeyPath + "/0"
)

// an hierarchical deterministic wallet
type HDWallet struct {
	id 			uuid.UUID
	walletType 	core.WalletType
	key 		*core.MasterDerivableKey // the node key from which all accounts are derived
	indexMapper map[string]uuid.UUID
	context 	*core.WalletContext
}

func NewHDWallet(key *core.MasterDerivableKey, context *core.WalletContext) *HDWallet {
	return &HDWallet{
		id:              uuid.New(),
		walletType:      core.HDWallet,
		key:        	 key,
		indexMapper: 	 make(map[string]uuid.UUID),
		context:		 context,
	}
}

// ID provides the ID for the wallet.
func (wallet *HDWallet) ID() uuid.UUID {
	return wallet.id
}

// Type provides the type of the wallet.
func (wallet *HDWallet) Type() core.WalletType {
	return wallet.walletType
}

// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
// This will error if an account with the name already exists.
func (wallet *HDWallet) CreateValidatorAccount(name string) (core.ValidatorAccount, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("account name is empty")
	}

	// Check if an account with the name already exists
	_, exists := wallet.indexMapper[name]
	if exists {
		return nil, fmt.Errorf("account %q already exists", name)
	}

	baseAccountPath := fmt.Sprintf(BaseAccountPath,len(wallet.indexMapper))
	// validator key
	validatorPath := fmt.Sprintf(ValidatorKeyPath,len(wallet.indexMapper))
	validatorKey,err := wallet.key.Derive(validatorPath)
	if err != nil {
		return nil, err
	}
	// withdrawal key
	withdrawalPath := fmt.Sprintf(WithdrawalKeyPath,len(wallet.indexMapper))
	withdrawalKey,err := wallet.key.Derive(withdrawalPath)
	if err != nil {
		return nil, err
	}

	// create ret account
	ret, err := NewValidatorAccount(
			name,
			validatorKey,
			withdrawalKey.PublicKey(),
			baseAccountPath,
			wallet.context,
		)
	if err != nil {
		return nil, err
	}

	// register new wallet and save portfolio
	reset := func() {
		delete(wallet.indexMapper,name)
	}
	wallet.indexMapper[name] = ret.ID()
	err = wallet.context.Storage.SaveAccount(ret)
	if err != nil {
		reset()
		return nil,err
	}
	err = wallet.context.Storage.SaveWallet(wallet)
	if err != nil {
		reset()
		return nil,err
	}

	return ret, nil
}

// Accounts provides all accounts in the wallet.
func (wallet *HDWallet) Accounts() <-chan core.ValidatorAccount {
	ch := make (chan core.ValidatorAccount,1024) // TODO - handle more? change from chan?
	go func() {
		for name := range wallet.indexMapper {
			id := wallet.indexMapper[name]
			account,err := wallet.AccountByID(id)
			if err != nil {
				continue
			}
			ch <- account
		}
		close(ch)
	}()

	return ch
}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByID(id uuid.UUID) (core.ValidatorAccount, error) {
	ret,err := wallet.context.Storage.OpenAccount(id)
	if err != nil {
		return nil,err
	}
	if ret == nil {
		return nil,nil
	}
	ret.SetContext(wallet.context)
	return ret,nil
}

func (wallet *HDWallet) SetContext(ctx *core.WalletContext) {
	wallet.context = ctx
}

// AccountByName provides a single account from the wallet given its name.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByName(name string) (core.ValidatorAccount, error) {
	id, exists := wallet.indexMapper[name]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}
	return wallet.AccountByID(id)
}