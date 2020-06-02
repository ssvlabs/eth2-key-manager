package core

import (
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// An account holds a key pair with the ability to do signatures and more
type Account interface {
	// ID provides the ID for the account.
	ID() uuid.UUID
	// Name provides the name for the account.
	Name() string
	// PublicKey provides the public key for the account.
	PublicKey() e2types.PublicKey
	// Path provides the path for the account.
	// Can be empty if the account is not derived from a path.
	Path() string
	// Sign signs data with the account.
	Sign(data []byte) (e2types.Signature, error)
	// lock will encrypt the seed, save it to memory and nil the plain text seed.
	// it will use an internally save locking password so it could be locked at all times
	Lock() error
	IsLocked() bool
	// unlock will decrypt the seed and save on memory
	// it needs a provided password
	Unlock(password []byte) error
}
