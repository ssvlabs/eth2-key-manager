package core

import (
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type AccountType = string
const (
	ValidatorAccount 	AccountType = "Validation"
	WithdrawalAccount	AccountType = "Withdrawal"
)

// An account holds a key pair with the ability to do signatures and more
type Account interface {
	// ID provides the ID for the account.
	ID() uuid.UUID
	// WalletID provides the ID for the wallet holding this account.
	WalletID() uuid.UUID
	// ID provides the ID for the account.
	Type() AccountType
	// Name provides the name for the account.
	Name() string
	// PublicKey provides the public key for the account.
	PublicKey() e2types.PublicKey
	// Path provides the path for the account.
	// Can be empty if the account is not derived from a path.
	Path() string
	// Sign signs data with the account.
	Sign(data []byte) (e2types.Signature,error)
	//
	SetContext(ctx *PortfolioContext)
}
