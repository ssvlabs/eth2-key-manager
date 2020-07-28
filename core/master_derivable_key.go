package core

import (
	"fmt"
	"github.com/google/uuid"
	util "github.com/wealdtech/go-eth2-util"
	"regexp"
)

const (
	BaseEIP2334Path = "m/12381/3600"
)

// MasterDerivableKey is responsible for managing privKey derivation, signing and  encryption.
// follows EIP 2333,2334
// MasterDerivableKey is not intended to be used a signing key, just as a medium for managing keys
type MasterDerivableKey struct {
	seed []byte
}

// base privKey is m / purpose / coin_type / as EIP 2334 defines
func MasterKeyFromSeed(seed []byte) (*MasterDerivableKey,error) {
	if seed == nil || len(seed) != 32 {
		return nil, fmt.Errorf("seed can't be nil or of length different than 32")
	}
	return &MasterDerivableKey{
		seed: seed,
	},nil
}

func (master *MasterDerivableKey) Derive(relativePath string) (*HDKey,error) {
	if !validateRelativePath(relativePath) {
		return nil, fmt.Errorf("invalid relative path. Example: /1/2/3")
	}

	//derive
	path := BaseEIP2334Path + relativePath
	key,err := util.PrivateKeyFromSeedAndPath(master.seed,path)
	if err != nil {
		return nil,err
	}

	// new id
	id := uuid.New()

	return &HDKey{
		id:      id,
		privKey:  key,
		path:    path,
	},nil
}

func validateRelativePath(relativePath string) bool {
	match, _ := regexp.MatchString(`^(\/(\d\d?\d?))+$`, relativePath)
	return match
}