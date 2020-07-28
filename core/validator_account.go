package core

import (
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

//type AccountType = string
//const (
//	ValidatorKey 	AccountType = "Validation"
//	WithdrawalKey	AccountType = "Withdrawal"
//)

// A validator account holds the information and actions needed by validator account keys.
// It holds 2 keys, a validation and a withdrawal key.
// As a minimum, the ValidatorAccount should have at least the validation key.
// Withdrawal key is not mandatory to be present.
type ValidatorAccount interface {
	// ID provides the ID for the account.
	ID() uuid.UUID
	// Name provides the name for the account.
	Name() string
	// ValidatorPublicKey provides the public key for the validation key.
	ValidatorPublicKey() e2types.PublicKey
	// WithdrawalPublicKey provides the public key for the withdrawal key.
	WithdrawalPublicKey() e2types.PublicKey
	// Sign signs data with the validation key.
	ValidationKeySign(data []byte) (e2types.Signature, error)
	//// Sign signs data with the withdrawal key.
	//WithdrawalKeySign(data []byte) (e2types.Signature,error)
	//
	SetContext(ctx *WalletContext)
}
