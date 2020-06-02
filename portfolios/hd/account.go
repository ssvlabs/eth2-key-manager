package hd

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type HDAccount struct {
	name string
	id uuid.UUID
	publicKey e2types.PublicKey
	secretKey e2types.PrivateKey
	path string
	lockPolicy core.LockablePolicy
}