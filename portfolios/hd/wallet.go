package hd

import (
	"fmt"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	types2 "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"github.com/wealdtech/go-indexer"
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

func NewHDWallet(name string, nodeKey *core.EncryptableSeed, path string, context *core.PortfolioContext) *HDWallet {
	return &HDWallet{
		name:            name,
		id:              uuid.New(),
		walletType:      core.HDWallet,
		nodeKey:         nodeKey,
		path:            path,
		accountsIndexer: indexer.New(),
		accountIds:      make([]uuid.UUID,0),
		context:		 context,
	}
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

// CreateAccount creates a new account in the wallet.
// This will error if an account with the name already exists.
// Will push to the new account the lock policy
func (wallet *HDWallet) CreateAccount(name string) (core.Account, error) {
	var retAccount *HDAccount
	if wallet.IsLocked() {
		return nil,fmt.Errorf("wallet is locked")
	}
	defer func() {
		if wallet.context.LockPolicy.LockAfterOperation(core.Creation) {
			wallet.Lock()
			if retAccount != nil {
				retAccount.Lock()
			}
		}
	}()

	// create account
	id := len(wallet.accountIds)
	path := fmt.Sprintf("%d",id)
	nodeBytes,err := wallet.deriveAccount(path,wallet.nodeKey.Seed())
	if err != nil {
		return nil,err
	}
	lockableKey := core.NewEncryptableSeed(nodeBytes.Marshal(),wallet.context.Encryptor)
	retAccount,err = newHDAccount(
		name,
		lockableKey,
		path,
		wallet.context,
	)

	// register new wallet and save portfolio
	wallet.accountIds = append(wallet.accountIds,retAccount.ID())
	wallet.accountsIndexer.Add(retAccount.ID(),name)
	err = wallet.context.Storage.SavePortfolio(wallet)
	if err != nil {
		wallet.accountsIndexer.Remove(retAccount.id,name)
		wallet.accountIds = wallet.accountIds[:len(wallet.accountIds)-1]
		return nil,err
	}
}

// Accounts provides all accounts in the wallet.
func (wallet *HDWallet) Accounts() <-chan core.Account {
	// TODO lockable policy
}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByID(id uuid.UUID) (core.Account, error) {
	// TODO lockable policy
}

// AccountByName provides a single account from the wallet given its name.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByName(name string) (core.Account, error) {
	// TODO lockable policy
}

func (wallet *HDWallet) Lock() error {
	return wallet.nodeKey.Encrypt(wallet.context.LockPassword)
}

func (wallet *HDWallet) IsLocked() bool {
	return wallet.nodeKey.IsEncrypted()
}

func (wallet *HDWallet) Unlock(password []byte) error {
	return wallet.nodeKey.Decrypt(password)
}

func (wallet *HDWallet) deriveAccount(path string, seed []byte) (*types2.BLSPrivateKey,error) {
	return types2.BLSPrivateKeyFromBytes(seed) // TODO - implement
}