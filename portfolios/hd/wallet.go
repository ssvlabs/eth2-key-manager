package hd

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	types2 "github.com/wealdtech/go-eth2-types/v2"
	"github.com/wealdtech/go-indexer"
)

// according to https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
const (
	WithdrawalKeyPath = "0"
	WithdrawalKeyName = "wallet_withdrawal_key_unique"
	ValidatorKeyPath = "0/%d"
)

// an hierarchical deterministic wallet
type HDWallet struct {
	name string
	id uuid.UUID
	walletType core.WalletType
	nodeKey *core.EncryptableSeed // the node key from which all accounts are derived
	path string
	accountsIndexer *indexer.Index  // maps indexs <> names
	accountIds []uuid.UUID
	context *core.PortfolioContext
}

func NewHDWallet(name string, nodeKey *core.EncryptableSeed, path string, context *core.PortfolioContext) (*HDWallet,error) {
	ret := &HDWallet{
		name:            name,
		id:              uuid.New(),
		walletType:      core.HDWallet,
		nodeKey:         nodeKey,
		path:            path,
		accountsIndexer: indexer.New(),
		accountIds:      make([]uuid.UUID,0),
		context:		 context,
	}

	_,err := ret.createKey(WithdrawalKeyName,WithdrawalKeyPath,core.WithdrawalAccount)
	if err != nil {
		return nil,err
	}

	return ret,nil
}

// ID provides the ID for the wallet.
func (wallet *HDWallet) ID() uuid.UUID {
	return wallet.id
}

// Name provides the name for the wallet.
func (wallet *HDWallet) Name() string {
	return wallet.name
}

// Type provides the type of the wallet.
func (wallet *HDWallet) Type() core.WalletType {
	return wallet.walletType
}

// GetWithdrawalAccount returns this wallet's withdrawal key pair in the wallet as described in EIP-2334.
// This will error if an account with the name already exists.
func (wallet *HDWallet) GetWithdrawalAccount() (core.Account, error) {
	return wallet.AccountByName(WithdrawalKeyName) // created when wallet is called with NewHDWallet
}

// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
// This will error if an account with the name already exists.
func (wallet *HDWallet) CreateValidatorAccount(name string) (core.Account, error) {
	path := fmt.Sprintf(ValidatorKeyPath,len(wallet.accountIds))
	return wallet.createKey(name,path,core.ValidatorAccount)
}

// Accounts provides all accounts in the wallet.
func (wallet *HDWallet) Accounts() <-chan core.Account {
}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByID(id uuid.UUID) (core.Account, error) {
	wallet.context.Storage.OpenAccount(id)
}

// AccountByName provides a single account from the wallet given its name.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByName(name string) (core.Account, error) {
	id,exists := wallet.accountsIndexer.ID(name)
	if !exists {
		return nil,fmt.Errorf("account not found")
	}

	return wallet.AccountByID(id)
}

func (wallet *HDWallet) deriveAccount(path string, seed []byte) (*types2.BLSPrivateKey,error) {
	return types2.BLSPrivateKeyFromBytes(seed) // TODO - implement
}

func (wallet *HDWallet) createKey(name string, path string, accountType core.AccountType) (core.Account, error) {
	var retAccount *HDAccount

	// create account
	nodeBytes,err := wallet.deriveAccount(path,wallet.nodeKey.Seed())
	if err != nil {
		return nil,err
	}
	lockableKey := core.NewEncryptableSeed(nodeBytes.Marshal(),wallet.context.Encryptor)
	retAccount,err = newHDAccount(
		name,
		accountType,
		lockableKey,
		path,
		wallet.context,
	)

	// register new wallet and save portfolio
	reset := func() {
		wallet.accountsIndexer.Remove(retAccount.id,name)
		wallet.accountIds = wallet.accountIds[:len(wallet.accountIds)-1]
	}
	wallet.accountIds = append(wallet.accountIds,retAccount.ID())
	wallet.accountsIndexer.Add(retAccount.ID(),name)
	err = wallet.context.Storage.SaveAccount(retAccount)
	if err != nil {
		reset()
		return nil,err
	}
	err = wallet.context.Storage.SaveWallet(wallet)
	if err != nil {
		reset()
		return nil,err
	}

	return retAccount,nil
}