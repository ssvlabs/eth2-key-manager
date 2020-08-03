package in_memory

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func (store *InMemStore) MarshalJSON() ([]byte, error) {
	var err error
	data := make(map[string]interface{})

	data["wallet"], err = json.Marshal(store.wallet)
	if err != nil {
		return nil, err
	}
	data["wallet"] = hex.EncodeToString(data["wallet"].([]byte))

	data["accounts"], err = json.Marshal(store.accounts)
	if err != nil {
		return nil, err
	}
	data["accounts"] = hex.EncodeToString(data["accounts"].([]byte))

	data["attMemory"], err = json.Marshal(store.attMemory)
	if err != nil {
		return nil, err
	}
	data["attMemory"] = hex.EncodeToString(data["attMemory"].([]byte))

	data["proposalMemory"], err = json.Marshal(store.proposalMemory)
	if err != nil {
		return nil, err
	}
	data["proposalMemory"] = hex.EncodeToString(data["proposalMemory"].([]byte))

	return json.Marshal(data)
}

func (store *InMemStore) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// wallet
	if val, exists := v["wallet"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		err = json.Unmarshal(byts, &store.wallet)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("could not find var: memory")
	}

	// accounts
	if val, exists := v["accounts"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		err = json.Unmarshal(byts, &store.accounts)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("could not find var: memory")
	}

	// attMemory
	if val, exists := v["attMemory"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		err = json.Unmarshal(byts, &store.attMemory)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("could not find var: attMemory")
	}

	// proposalMemory
	if val, exists := v["proposalMemory"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		err = json.Unmarshal(byts, &store.proposalMemory)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("could not find var: proposalMemory")
	}

	return nil
}
