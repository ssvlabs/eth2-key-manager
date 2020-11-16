package core

import (
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// HDKey is a derived key from MasterDerivableKey which is able to sign messages, return thee public key and more.
type HDKey struct {
	id      uuid.UUID
	privKey e2types.PrivateKey
	path    string
}

// NewHDKeyFromPrivateKey is the constructor of HDKey
func NewHDKeyFromPrivateKey(priv []byte, path string) (*HDKey, error) {
	key, err := e2types.BLSPrivateKeyFromBytes(priv)
	if err != nil {
		return nil, err
	}

	return &HDKey{
		id:      uuid.New(),
		privKey: key,
		path:    path,
	}, nil
}

func (key *HDKey) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = key.id
	data["privKey"] = hex.EncodeToString(key.privKey.Marshal())
	data["path"] = key.path

	return json.Marshal(data)
}

func (key *HDKey) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if val, exists := v["id"]; exists {
		var err error
		if key.id, err = uuid.Parse(val.(string)); err != nil {
			return err
		}
	} else {
		return errors.New("could not find var: id")
	}

	if val, exists := v["path"]; exists {
		key.path = val.(string)
	} else {
		return errors.New("could not find var: path")
	}

	if val, exists := v["privKey"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}

		if key.privKey, err = e2types.BLSPrivateKeyFromBytes(byts); err != nil {
			return err
		}
	} else {
		return errors.New("could not find var: privKey")
	}

	return nil
}

func (key *HDKey) PublicKey() e2types.PublicKey {
	return key.privKey.PublicKey()
}

func (key *HDKey) Sign(data []byte) (e2types.Signature, error) {
	return key.privKey.Sign(data), nil
}

func (key *HDKey) Path() string {
	return key.path
}
