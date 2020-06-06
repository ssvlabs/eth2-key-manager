package wallet_hd

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type HDAccount struct {
	name string
	id uuid.UUID
	accountType core.AccountType
	publicKey e2types.PublicKey
	key *core.DerivableKey
	context *core.PortfolioContext
}

func newHDAccount(name string,
	accountType core.AccountType,
	key *core.DerivableKey,
	context *core.PortfolioContext) (*HDAccount,error) {
	return &HDAccount{
		name:         name,
		id:           uuid.New(),
		accountType:  accountType,
		publicKey:    key.Key.PublicKey(),
		key:    	  key,
		context:	  context,
	},nil
}

// ID provides the ID for the account.
func (account *HDAccount) ID() uuid.UUID {
	return account.id
}

// WalletID provides the ID for the wallet holding this account.
func (account *HDAccount) WalletID() uuid.UUID {
	return account.context.WalletId
}

// ID provides the ID for the account.
func (account *HDAccount) Type() core.AccountType {
	return account.accountType
}

// Name provides the name for the account.
func (account *HDAccount) Name() string {
	return account.name
}

// PublicKey provides the public key for the account.
func (account *HDAccount) PublicKey() e2types.PublicKey {
	return account.publicKey
}

// Path provides the path for the account.
// Can be empty if the account is not derived from a path.
func (account *HDAccount) Path() string {
	return account.key.Path
}

// Sign signs data with the account.
func (account *HDAccount) Sign(data []byte) e2types.Signature {
	return account.key.Key.Sign(data)
}
