package core

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	util "github.com/wealdtech/go-eth2-util"
)

// EIP2334 paths.
const (
	// BaseEIP2334Path is the base EIP2334 path.
	BaseEIP2334Path = "m/12381/3600"
)

// MasterDerivableKey is responsible for managing privKey derivation, signing and  encryption.
// follows EIP 2333,2334
// MasterDerivableKey is not intended to be used a signing key, just as a medium for managing keys
type MasterDerivableKey struct {
	seed    []byte
	network Network
}

// MasterKeyFromSeed is the constructor of MasterDerivableKey.
// Base privKey is m / purpose / coin_type / as EIP 2334 defines
func MasterKeyFromSeed(seed []byte, network Network) (*MasterDerivableKey, error) {
	if seed == nil || len(seed) == 0 {
		return nil, errors.New("seed can't be nil or length 0")
	}
	return &MasterDerivableKey{
		seed:    seed,
		network: network,
	}, nil
}

// Derive derives a HD key based on the given relative path.
func (master *MasterDerivableKey) Derive(relativePath string) (*HDKey, error) {
	if !validateRelativePath(relativePath) {
		return nil, errors.New("invalid relative path. Example: /1/2/3")
	}

	path := master.network.FullPath(relativePath)
	key, err := util.PrivateKeyFromSeedAndPath(master.seed, path)
	if err != nil {
		return nil, err
	}

	return &HDKey{
		id:      uuid.New(),
		privKey: key,
		path:    path,
	}, nil
}

func validateRelativePath(relativePath string) bool {
	match, _ := regexp.MatchString(`^(\/(\d\d?\d?))+$`, relativePath)
	return match
}
