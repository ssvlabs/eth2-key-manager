package hd

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type HDAccount struct {
	name string
	id uuid.UUID
	publicKey e2types.PublicKey
	secretKey *core.EncryptableSeed
	path string
	lockPolicy core.LockablePolicy
	lockPassword []byte // only used internally for quick lock
}

func newHDAccount(name string, secretKey *core.EncryptableSeed, path string, lockPolicy core.LockablePolicy, lockPassword []byte) (*HDAccount,error) {
	if secretKey.IsEncrypted() {
		return nil,fmt.Errorf("account is locked")
	}

	priv,err := e2types.BLSPrivateKeyFromBytes(secretKey.Seed())
	if err != nil {
		return nil,err
	}

	return &HDAccount{
		name:         name,
		id:           uuid.New(),
		publicKey:    priv.PublicKey(),
		secretKey:    secretKey,
		path:         path,
		lockPolicy:   lockPolicy,
		lockPassword: lockPassword,
	},nil
}

// ID provides the ID for the account.
func (account *HDAccount) ID() uuid.UUID {
	return account.id
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
	return account.path
}

// Sign signs data with the account.
func (account *HDAccount) Sign(data []byte) (e2types.Signature, error) {

}

func (account *HDAccount) Lock() error {
	return account.secretKey.Encrypt(account.lockPassword)
}

func (account *HDAccount) IsLocked() bool {
	return account.secretKey.IsEncrypted()
}

func (account *HDAccount) Unlock(password []byte) error {
	return account.secretKey.Decrypt(password)
}