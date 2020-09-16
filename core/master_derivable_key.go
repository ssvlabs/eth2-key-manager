package core

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	util "github.com/wealdtech/go-eth2-util"
)

// EIP2334 paths.
const (
	// BaseEIP2334Path is the base EIP2334 of the Main Net.
	BaseEIP2334Path = "m/12381/3600"

	// BaseTestEIP2334Path is the base EIP2334 of the Test Net.
	BaseTestEIP2334Path = "m/12381/3599"

	// BaseLaunchTestEIP2334Path is the base EIP2334 of the Launch Test Net.
	BaseLaunchTestEIP2334Path = "m/12381/3598"
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
		return nil, fmt.Errorf("seed can't be nil or length 0")
	}
	return &MasterDerivableKey{
		seed:    seed,
		network: network,
	}, nil
}

// Derive derives a HD key based on the given relative path.
func (master *MasterDerivableKey) Derive(relativePath string) (*HDKey, error) {
	if !validateRelativePath(relativePath) {
		return nil, fmt.Errorf("invalid relative path. Example: /1/2/3")
	}

	// Derive key
	path := master.network.FullPath(relativePath)
	key, err := util.PrivateKeyFromSeedAndPath(master.seed, path)
	if err != nil {
		return nil, err
	}

	// Create key ID
	id := uuid.New()

	return &HDKey{
		id:      id,
		privKey: key,
		path:    path,
	}, nil
}

func validateRelativePath(relativePath string) bool {
	match, _ := regexp.MatchString(`^(\/(\d\d?\d?))+$`, relativePath)
	return match
}
