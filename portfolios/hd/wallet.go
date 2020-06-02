package hd

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-types/v2"
	"github.com/wealdtech/go-indexer"
)

// an hierarchical deterministic wallet
type HDWallet struct {
	name string
	id uuid.UUID
	walletType core.WalletType
	nodeKey types.BLSPrivateKey // the node key from which all accounts are derived
	path string
	lockPolicy core.LockablePolicy
	accountsIndexer indexer.Index  // maps indexs <> names
	accountIds []uuid.UUID
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
