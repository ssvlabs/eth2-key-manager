package nd

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// MarshalJSON is the custom JSON marshaler
func (wallet *Wallet) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = wallet.id
	data["type"] = wallet.walletType
	data["indexMapper"] = wallet.indexMapper

	return json.Marshal(data)
}

// UnmarshalJSON is the custom JSON unmarshaler
func (wallet *Wallet) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	var err error

	// id
	if val, exists := v["id"]; exists {
		wallet.id, err = uuid.Parse(val.(string))
		if err != nil {
			return err
		}
	} else {
		return errors.New("could not find var: id")
	}

	// type
	if val, exists := v["type"]; exists {
		wallet.walletType = val.(string)
	} else {
		return errors.New("could not find var: type")
	}

	// indexMapper
	if val, exists := v["indexMapper"]; exists {
		wallet.indexMapper = make(map[string]uuid.UUID)
		for k, v := range val.(map[string]interface{}) {
			wallet.indexMapper[k], err = uuid.Parse(v.(string))
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("could not find var: indexMapper")
	}

	return nil
}
