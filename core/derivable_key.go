package core

import (
	"fmt"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	"regexp"
)

const (
	BaseEIP2334Path = "m/12381/3600"
)

// follows EIP 2333,2334
type DerivableKey struct {
	seed []byte
	Key  *e2types.BLSPrivateKey
	Path string
}

// base key is m / purpose / coin_type / as EIP 2334 defines
func BaseKeyFromSeed(seed []byte) (*DerivableKey,error) {
	key,err := util.PrivateKeyFromSeedAndPath(seed,BaseEIP2334Path)
	if err != nil {
		return nil,err
	}

	return &DerivableKey{seed:seed,Key:key,Path:BaseEIP2334Path},nil
}

func (baseKey *DerivableKey) Derive(relativePath string) (*DerivableKey,error) {
	if !validateRelativePath(relativePath) {
		return nil, fmt.Errorf("invalid relative path. Example: /1/2/3")
	}

	path := baseKey.Path + relativePath
	key,err := util.PrivateKeyFromSeedAndPath(baseKey.seed,path)
	if err != nil {
		return nil,err
	}

	return &DerivableKey{seed:baseKey.seed,Key:key,Path:path},nil
}

func validateRelativePath(relativePath string) bool {
	match, _ := regexp.MatchString(`^(\/(\d\d?\d?))+$`, relativePath)
	return match
}