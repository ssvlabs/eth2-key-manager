package core

import (
	"encoding/hex"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type SimpleAccount struct {
	privateKey *e2types.BLSPrivateKey
	id uuid.UUID
}


// a simple account that's locally generates it's own private key
func NewSimpleAccount() *SimpleAccount {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil
	}

	priv,err := e2types.GenerateBLSPrivateKey()
	if err != nil {
		return nil
	}

	return &SimpleAccount{privateKey:priv,id:uuid.New()}
}


// WalletID provides the ID for the wallet holding this account.
func (account *SimpleAccount) WalletID() uuid.UUID {
	return uuid.New()
}

// ID provides the ID for the account.
func (account *SimpleAccount) ID() uuid.UUID {
	return account.id
}

// Name provides the name for the account.
func (account *SimpleAccount) Name() string {
	return hex.EncodeToString(account.privateKey.PublicKey().Marshal())
}

// ID provides the ID for the account.
func (account *SimpleAccount) Type() AccountType {
	return ValidatorAccount
}

// PublicKey provides the public key for the account.
func (account *SimpleAccount) PublicKey() e2types.PublicKey {
	return account.privateKey.PublicKey()
}

// Path provides the path for the account.
// Can be empty if the account is not derived from a path.
func (account *SimpleAccount) Path() string {
	return "m"
}

// Sign signs data with the account.
func (account *SimpleAccount) Sign(data []byte) e2types.Signature {
	return account.privateKey.Sign(data)
}
