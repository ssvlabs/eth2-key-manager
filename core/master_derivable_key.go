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
// the private key or the seed are never held in memory but rather fetched ad-hoc from Storage.
// how secure the Storage is for saving the seed (including encryption) is on the Storage implementation
// follows EIP 2333,2334
//
// MasterDerivableKey is not intended to be used a signing key, just as a medium for managing keys
type MasterDerivableKey struct {
	Storage Storage
}

// base privKey is m / purpose / coin_type / as EIP 2334 defines
func MasterKeyFromSeed(storage Storage) (*MasterDerivableKey,error) {
	return &MasterDerivableKey{
		Storage: storage,
	},nil
}

func (baseKey *MasterDerivableKey) Derive(relativePath string) (*HDKey,error) {
	if !validateRelativePath(relativePath) {
		return nil, fmt.Errorf("invalid relative path. Example: /1/2/3")
	}

	// fetch priv key
	seed,err := baseKey.tempFetchSeed()
	if err != nil {
		return nil,err
	}
	if seed == nil {
		return nil,fmt.Errorf("seed is nil")
	}

	//derive
	path := BaseEIP2334Path + relativePath
	key,err := util.PrivateKeyFromSeedAndPath(seed,path)
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

// TODO - this is a limitation of the util library that we use as it can't derive relative path but only absolute from the seed.
func (key *MasterDerivableKey) tempFetchSeed() ([]byte,error) {
	return key.Storage.SecurelyFetchPortfolioSeed()
}

func validateRelativePath(relativePath string) bool {
	match, _ := regexp.MatchString(`^(\/(\d\d?\d?))+$`, relativePath)
	return match
}