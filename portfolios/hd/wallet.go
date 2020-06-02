package hd

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	util "github.com/wealdtech/go-eth2-util"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"github.com/wealdtech/go-indexer"
)

// an hierarchical deterministic wallet
type HDWallet struct {
	name string
	id uuid.UUID
	encryptor types.Encryptor
	walletType core.WalletType
	nodeKey *core.EncryptableSeed // the node key from which all accounts are derived
	path string
	lockPolicy core.LockablePolicy
	accountsIndexer *indexer.Index  // maps indexs <> names
	accountIds []uuid.UUID
	lockPassword []byte // only used internally for quick lock
}

func NewHDWallet(name string, nodeKey *core.EncryptableSeed, path string, lockPolicy core.LockablePolicy, encryptor types.Encryptor, lockPassword []byte) *HDWallet {
	return &HDWallet{
		name:            name,
		id:              uuid.New(),
		encryptor:		 encryptor,
		walletType:      core.HDWallet,
		nodeKey:         nodeKey,
		path:            path,
		lockPolicy:      lockPolicy,
		accountsIndexer: indexer.New(),
		accountIds:      make([]uuid.UUID,0),
		lockPassword:	lockPassword,
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
		if wallet.lockPolicy.LockAfterOperation(core.Creation) {
			wallet.Lock()
			if retAccount != nil {
				retAccount.Lock()
			}
		}
	}()

	// create account
	id := len(wallet.accountIds)
	path := fmt.Sprintf("%d",id)
	nodeBytes,err := util.PrivateKeyFromSeedAndPath(wallet.nodeKey.Seed(),path) // TODO - this will not work as we do not give an 'm' in the path
	if err != nil {
		return nil,err
	}
	lockableKey := core.NewEncryptableSeed(nodeBytes.Marshal(),wallet.encryptor)
	retAccount,err = newHDAccount(
		name,
		lockableKey,
		path,
		wallet.lockPolicy,
		wallet.lockPassword,
	)

	// register new wallet and save portfolio
	wallet.accountIds = append(wallet.accountIds,retAccount.ID())
	wallet.accountsIndexer.Add(retAccount.ID(),name)
	err = wallet.storage.SavePortfolio(wallet)
	if err != nil {
		wallet.walletsIndexer.Remove(retWallet.id,name)
		wallet.walletIds = portfolio.walletIds[:len(wallet.walletIds)-1]
		return nil,err
	}
}

// Accounts provides all accounts in the wallet.
func (wallet *HDWallet) Accounts() <-chan core.Account {

}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByID(id uuid.UUID) (core.Account, error) {

}

// AccountByName provides a single account from the wallet given its name.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByName(name string) (core.Account, error) {

}

func (wallet *HDWallet) Lock() error {
	return wallet.nodeKey.Encrypt(wallet.lockPassword)
}

func (wallet *HDWallet) IsLocked() bool {
	return wallet.nodeKey.IsEncrypted()
}

func (wallet *HDWallet) Unlock(password []byte) error {
	return wallet.nodeKey.Decrypt(password)
}