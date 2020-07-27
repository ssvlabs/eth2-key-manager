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
	validationKey *core.HDKey
	withdrawalKey *core.HDKey
	context *core.WalletContext
}

func (account *HDAccount) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = account.id
	data["name"] = account.name
	data["validationKey"] = account.validationKey
	data["withdrawalKey"] = account.withdrawalKey
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

	// validation key
	if val, exists := v["validationKey"]; exists {
		byts,err := json.Marshal(val)
		if err != nil {
			return err
		}
		key := &core.HDKey{}
		err = json.Unmarshal(byts,key)
		if err != nil {
			return err
		}
		account.validationKey = key
	} else {return fmt.Errorf("could not find var: key")}

	// withdrawal Key
	if val, exists := v["withdrawalKey"]; exists {
		byts,err := json.Marshal(val)
		if err != nil {
			return err
		}
		key := &core.HDKey{}
		err = json.Unmarshal(byts,key)
		if err != nil {
			return err
		}
		account.withdrawalKey = key
	} else {return fmt.Errorf("could not find var: key")}

	return nil
}

func NewValidatorAccount(name string,
	validationKey *core.HDKey,
	withdrawalKey *core.HDKey,
	context *core.WalletContext) (*HDAccount,error) {
	return &HDAccount{
		name:         	 name,
		id:          	 uuid.New(),
		validationKey:	 validationKey,
		withdrawalKey:	 withdrawalKey,
		context:	  	 context,
	},nil
}

// ID provides the ID for the account.
func (account *HDAccount) ID() uuid.UUID {
	return account.id
}

// Name provides the name for the account.
func (account *HDAccount) Name() string {
	return account.name
}

// ValidatorPublicKey provides the public key for the account.
func (account *HDAccount) ValidatorPublicKey() e2types.PublicKey {
	return account.validationKey.PublicKey()
}

// WithdrawalPublicKey provides the public key for the account.
func (account *HDAccount) WithdrawalPublicKey() e2types.PublicKey {
	return account.withdrawalKey.PublicKey()
}

// Sign signs data with the account.
func (account *HDAccount) ValidationKeySign(data []byte) (e2types.Signature,error) {
	return account.validationKey.Sign(data)
}

// Sign signs data with the withdrawal key.
func (account *HDAccount) WithdrawalKeySign(data []byte) (e2types.Signature,error) {
	if account.withdrawalKey == nil {
		return nil, fmt.Errorf("withdrawal key not present")
	}
	return account.withdrawalKey.Sign(data)
}


func (account *HDAccount) SetContext(ctx *core.WalletContext) {
	account.context = ctx
}