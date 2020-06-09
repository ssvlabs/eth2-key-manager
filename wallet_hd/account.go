package wallet_hd

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type HDAccount struct {
	name string
	id uuid.UUID
	accountType core.AccountType
	key *core.DerivableKey
	context *core.PortfolioContext
}

func (account *HDAccount) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = account.id
	data["name"] = account.name
	data["type"] = account.accountType
	data["key"] = account.key
	return json.Marshal(data)
}

func (account *HDAccount) UnmarshalJSON(data []byte) error {
	// parse
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	var err error

	// id
	if val, exists := v["id"]; exists {
		account.id,err = uuid.Parse(val.(string))
		if err != nil {
			return err
		}
	} else {return fmt.Errorf("could not find var: id")}

	// name
	if val, exists := v["name"]; exists {
		account.name = val.(string)
	} else {return fmt.Errorf("could not find var: id")}

	// type
	if val, exists := v["type"]; exists {
		account.accountType = val.(string)
	} else {return fmt.Errorf("could not find var: id")}

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
		account.key = key
	} else {return fmt.Errorf("could not find var: key")}

	return nil
}

func newHDAccount(name string,
	accountType core.AccountType,
	key *core.DerivableKey,
	context *core.PortfolioContext) (*HDAccount,error) {
	return &HDAccount{
		name:         name,
		id:           uuid.New(),
		accountType:  accountType,
		key:    	  key,
		context:	  context,
	},nil
}

// ID provides the ID for the account.
func (account *HDAccount) ID() uuid.UUID {
	return account.id
}

// WalletID provides the ID for the wallet holding this account.
func (account *HDAccount) WalletID() uuid.UUID {
	return account.context.WalletId
}

// ID provides the ID for the account.
func (account *HDAccount) Type() core.AccountType {
	return account.accountType
}

// Name provides the name for the account.
func (account *HDAccount) Name() string {
	return account.name
}

// PublicKey provides the public key for the account.
func (account *HDAccount) PublicKey() e2types.PublicKey {
	return account.key.PublicKey()
}

// Path provides the path for the account.
// Can be empty if the account is not derived from a path.
func (account *HDAccount) Path() string {
	return account.key.GetPath()
}

// Sign signs data with the account.
func (account *HDAccount) Sign(data []byte) e2types.Signature {
	return account.key.Sign(data)
}
