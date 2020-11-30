package in_memory

import (
	"encoding/hex"
	"encoding/json"

	"github.com/bloxapp/eth2-key-manager/wallets/nd"

	hd2 "github.com/bloxapp/eth2-key-manager/wallets/hd"

	"github.com/pkg/errors"

	"github.com/bloxapp/eth2-key-manager/core"
)

func (store *InMemStore) MarshalJSON() ([]byte, error) {
	var err error
	data := make(map[string]interface{})

	data["network"] = hex.EncodeToString([]byte(store.network))

	data["wallet"], err = json.Marshal(store.wallet)
	if err != nil {
		return nil, err
	}
	data["wallet"] = hex.EncodeToString(data["wallet"].([]byte))

	data["walletType"] = store.wallet.Type()

	data["accounts"], err = json.Marshal(store.accounts)
	if err != nil {
		return nil, err
	}
	data["accounts"] = hex.EncodeToString(data["accounts"].([]byte))

	data["highestAtt"], err = json.Marshal(store.highestAttestation)
	if err != nil {
		return nil, err
	}
	data["highestAtt"] = hex.EncodeToString(data["highestAtt"].([]byte))

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

	// network
	if val, exists := v["network"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}

		store.network = core.NetworkFromString(string(byts))
	} else {
		return errors.New("could not find var: network")
	}

	// wallet
	if walletType, exists := v["walletType"]; exists {
		if val, exists := v["wallet"]; exists {
			byts, err := hex.DecodeString(val.(string))
			if err != nil {
				return err
			}

			if walletType == core.HDWallet {
				hd := &hd2.HDWallet{}
				err = json.Unmarshal(byts, &hd)
				if err != nil {
					return err
				}
				store.wallet = hd
			} else if walletType == core.NDWallet {
				nd := &nd.NDWallet{}
				err = json.Unmarshal(byts, &nd)
				if err != nil {
					return err
				}
				store.wallet = nd
			} else {
				return errors.Errorf("unknown wallet type %s", walletType)
			}
		} else {
			return errors.New("could not find var: wallet")
		}
	} else {
		return errors.New("could not find var: walletType")
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
		return errors.New("could not find var: accounts")
	}

	// highest att.
	if val, exists := v["highestAtt"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		err = json.Unmarshal(byts, &store.highestAttestation)
		if err != nil {
			return err
		}
	} else {
		return errors.New("could not find var: highestAtt")
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
		return errors.New("could not find var: proposalMemory")
	}

	return nil
}
