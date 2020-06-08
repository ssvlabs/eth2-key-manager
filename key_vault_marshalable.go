package KeyVault

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
)

func (vault *KeyVault) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = vault.id
	data["enableSimpleSigner"] = vault.enableSimpleSigner
	data["indexMapper"] = vault.indexMapper
	data["key"] = vault.key
	return json.Marshal(data)
}

func (vault *KeyVault) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	var err error

	// id
	if val, exists := v["id"]; exists {
		vault.id,err = uuid.Parse(val.(string))
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: id")}

	// simple signer
	if val, exists := v["enableSimpleSigner"]; exists {
		vault.enableSimpleSigner = val.(bool)
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: enableSimpleSigner")}

	// indexMapper
	if val, exists := v["indexMapper"]; exists {
		vault.indexMapper = make(map[string]uuid.UUID)
		for k,v := range val.(map[string]interface{}) {
			vault.indexMapper[k],err = uuid.Parse(v.(string))
			if err != nil {
				return err
			}
		}
	} else {return fmt.Errorf("could not find var: indexMapper")}

	// key
	if val, exists := v["key"]; exists {
		byts,err := json.Marshal(val)
		if err != nil {
			return err
		}
		key := &core.DerivableKey{}
		err = json.Unmarshal(byts,key)
		if err != nil {
			return err
		}
		vault.key = key
	} else {return fmt.Errorf("could not find var: key")}

	return nil
}