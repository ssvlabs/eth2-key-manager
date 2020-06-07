package wallet_hd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

func (wallet *HDWallet) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = wallet.id
	data["name"] = wallet.name
	data["type"] = wallet.walletType
	data["indexMapper"] = wallet.indexMapper

	return json.Marshal(data)
}

func (wallet *HDWallet) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	var err error

	// id
	if val, exists := v["id"]; exists {
		wallet.id,err = uuid.Parse(val.(string))
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: id")}

	// name
	if val, exists := v["name"]; exists {
		wallet.name = val.(string)
	} else {return fmt.Errorf("could not find var: id")}

	// type
	if val, exists := v["type"]; exists {
		wallet.walletType = val.(string)
	} else {return fmt.Errorf("could not find var: id")}

	// indexMapper
	if val, exists := v["indexMapper"]; exists {
		wallet.indexMapper = make(map[string]uuid.UUID)
		for k,v := range val.(map[string]interface{}) {
			wallet.indexMapper[k],err = uuid.Parse(v.(string))
			if err != nil {
				return err
			}
		}
	} else {return fmt.Errorf("could not find var: indexMapper")}
	return nil
}