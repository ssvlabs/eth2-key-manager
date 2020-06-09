package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	"regexp"
	"strings"
)

const (
	BaseEIP2334Path = "m/12381/3600"
)

// DerivableKey is responsible for managing key derivation, signing and  encryption.
// follows EIP 2333,2334
type DerivableKey struct {
	seed []byte
	key  *e2types.BLSPrivateKey
	path string
}
func (key *DerivableKey) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["seed"] = hex.EncodeToString(key.seed)
	data["path"] = key.path

	return json.Marshal(data)
}

func (key *DerivableKey) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// path
	if val, exists := v["path"]; exists {
		key.path = val.(string)
	} else {return fmt.Errorf("could not find var: id")}

	// seed
	if val, exists := v["seed"]; exists {
		byts,err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		key.seed = byts

		baseKey,err := BaseKeyFromSeed(byts)
		if err != nil {
			return err
		}
		relativePath := strings.Replace(key.path,BaseEIP2334Path,"",1)
		if len(relativePath) > 0 {
			derivedKey,err := baseKey.Derive(relativePath)
			if err != nil {
				return err
			}
			key.key = derivedKey.key
		} else {
			key.key = baseKey.key
		}

	} else {return fmt.Errorf("could not find var: id")}

	return nil
}

// base key is m / purpose / coin_type / as EIP 2334 defines
func BaseKeyFromSeed(seed []byte) (*DerivableKey,error) {
	key,err := util.PrivateKeyFromSeedAndPath(seed,BaseEIP2334Path)
	if err != nil {
		return nil,err
	}

	return &DerivableKey{seed:seed,key:key,path:BaseEIP2334Path},nil
}

func (baseKey *DerivableKey) Derive(relativePath string) (*DerivableKey,error) {
	if !validateRelativePath(relativePath) {
		return nil, fmt.Errorf("invalid relative path. Example: /1/2/3")
	}

	path := baseKey.path + relativePath
	key,err := util.PrivateKeyFromSeedAndPath(baseKey.seed,path)
	if err != nil {
		return nil,err
	}

	return &DerivableKey{seed:baseKey.seed,key:key,path:path},nil
}

func (baseKey *DerivableKey) PublicKey() e2types.PublicKey {
	return baseKey.key.PublicKey()
}

func (baseKey *DerivableKey) Sign(data []byte) e2types.Signature {
	return baseKey.key.Sign(data)
}

func (baseKey *DerivableKey) GetPath() string {
	return baseKey.path
}

func validateRelativePath(relativePath string) bool {
	match, _ := regexp.MatchString(`^(\/(\d\d?\d?))+$`, relativePath)
	return match
}