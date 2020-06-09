package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	"regexp"
)

const (
	BaseEIP2334Path = "m/12381/3600"
)

// DerivableKey is responsible for managing privKey derivation, signing and  encryption.
// the private key or the seed are never held in memory but rather fetched ad-hoc from Storage.
// how secure the Storage is for saving the seed (including encryption) is on the Storage implementation
// follows EIP 2333,2334
type DerivableKey struct {
	id      uuid.UUID
	pubKey  e2types.PublicKey
	path    string
	Storage Storage
}
func (key *DerivableKey) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = key.id
	data["pubkey"] = hex.EncodeToString(key.pubKey.Marshal())
	data["path"] = key.path

	return json.Marshal(data)
}

func (key *DerivableKey) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// id
	if val, exists := v["id"]; exists {
		var err error
		key.id,err = uuid.Parse(val.(string))
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: id")}

	// path
	if val, exists := v["path"]; exists {
		key.path = val.(string)
	} else {return fmt.Errorf("could not find var: id")}


	// pubkey
	if val, exists := v["pubkey"]; exists {
		byts,err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		key.pubKey,err = e2types.BLSPublicKeyFromBytes(byts)
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: id")}

	return nil
}

// base privKey is m / purpose / coin_type / as EIP 2334 defines
func BaseKeyFromSeed(seed []byte, storage Storage) (*DerivableKey,error) {
	key,err := util.PrivateKeyFromSeedAndPath(seed,BaseEIP2334Path)
	if err != nil {
		return nil,err
	}

	id := uuid.New()

	return &DerivableKey{
		Storage: storage,
		id:      id,
		pubKey:  key.PublicKey(),
		path:    BaseEIP2334Path,
	},nil
}

func (baseKey *DerivableKey) Derive(relativePath string) (*DerivableKey,error) {
	if !validateRelativePath(relativePath) {
		return nil, fmt.Errorf("invalid relative path. Example: /1/2/3")
	}

	// fetch priv key
	seed,err := baseKey.tempFetchSeed()
	if err != nil {
		return nil,err
	}

	//derive
	path := baseKey.path + relativePath
	key,err := util.PrivateKeyFromSeedAndPath(seed,path)
	if err != nil {
		return nil,err
	}

	// new id
	id := uuid.New()

	return &DerivableKey{
		Storage: baseKey.Storage,
		id:      id,
		pubKey:  key.PublicKey(),
		path:    path,
	},nil
}

func (key *DerivableKey) PublicKey() e2types.PublicKey {
	return key.pubKey
}

func (key *DerivableKey) Sign(data []byte) (e2types.Signature,error) {
	privKey,err := key.tempFetchPrivKey()
	if err != nil {
		return nil,err
	}
	return privKey.Sign(data),nil
}

func (key *DerivableKey) Path() string {
	return key.path
}

func (key *DerivableKey) tempFetchPrivKey() (e2types.PrivateKey,error) {
	seed,err := key.tempFetchSeed()
	if err != nil {
		return nil,err
	}

	return util.PrivateKeyFromSeedAndPath(seed,key.Path())
}

// TODO - this is a limitation of the util library that we use as it can't derive relative path but only absolute from the seed.
func (key *DerivableKey) tempFetchSeed() ([]byte,error) {
	return key.Storage.SecurelyFetchPortfolioSeed()
}

func validateRelativePath(relativePath string) bool {
	match, _ := regexp.MatchString(`^(\/(\d\d?\d?))+$`, relativePath)
	return match
}