package inmemory

import (
	"encoding/hex"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ssvlabs/eth2-key-manager/core"
	hd2 "github.com/ssvlabs/eth2-key-manager/wallets/hd"
	"github.com/ssvlabs/eth2-key-manager/wallets/nd"
)

// MarshalJSON is the custom JSON marshaler
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

	data["highestProposal"], err = json.Marshal(store.highestProposal)
	if err != nil {
		return nil, err
	}
	data["highestProposal"] = hex.EncodeToString(data["highestProposal"].([]byte))

	return json.Marshal(data)
}

// UnmarshalJSON is the custom JSON unmarshaler
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

		store.network, err = core.NetworkFromString(string(byts))
		if err != nil {
			return err
		}
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

			switch walletType {
			case core.HDWallet:
				hd := &hd2.Wallet{}
				err = json.Unmarshal(byts, &hd)
				if err != nil {
					return err
				}
				store.wallet = hd
			case core.NDWallet:
				nd := &nd.Wallet{}
				err = json.Unmarshal(byts, &nd)
				if err != nil {
					return err
				}
				store.wallet = nd
			default:
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

	// highestProposal
	if val, exists := v["highestProposal"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		err = json.Unmarshal(byts, &store.highestProposal)
		if err != nil {
			return err
		}
	} else {
		return errors.New("could not find var: highestProposal")
	}

	return nil
}
