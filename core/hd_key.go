package core

import (
	"encoding/hex"
	"encoding/json"

	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// HDKey is a derived key from MasterDerivableKey which is able to sign messages, return thee public key and more.
type HDKey struct {
	id      uuid.UUID
	privKey *bls.SecretKey
	path    string
}

// NewHDKeyFromPrivateKey is the constructor of HDKey
func NewHDKeyFromPrivateKey(priv []byte, path string) (*HDKey, error) {
	sk := &bls.SecretKey{}
	if err := sk.SetHexString(hex.EncodeToString(priv)); err != nil {
		return nil, err
	}

	return &HDKey{
		id:      uuid.New(),
		privKey: sk,
		path:    path,
	}, nil
}

// MarshalJSON is the custom JSON marshaler
func (key *HDKey) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = key.id
	data["privKey"] = hex.EncodeToString(key.privKey.Serialize())
	data["path"] = key.path

	return json.Marshal(data)
}

// UnmarshalJSON is the custom JSON unmarshaler
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
		key.privKey = &bls.SecretKey{}

		if err := key.privKey.SetHexString(val.(string)); err != nil {
			return err
		}
	} else {
		return errors.New("could not find var: privKey")
	}

	return nil
}

// PublicKey returns the public key
func (key *HDKey) PublicKey() *bls.PublicKey {
	return key.privKey.GetPublicKey()
}

// Sign signs the given data
func (key *HDKey) Sign(data []byte) ([]byte, error) {
	return key.privKey.SignByte(data).Serialize(), nil
}

// Path returns path
func (key *HDKey) Path() string {
	return key.path
}
