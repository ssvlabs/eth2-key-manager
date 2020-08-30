package wallet_hd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth-key-manager/core"
	"github.com/bloxapp/eth-key-manager/eth1_deposit"
)

type HDAccount struct {
	name string
	// holds the base path from which the account was dervied
	// for eip2334 should be m/12381/3600/<index>
	basePath         string
	id               uuid.UUID
	validationKey    *core.HDKey
	withdrawalPubKey e2types.PublicKey
	context          *core.WalletContext
}

func (account *HDAccount) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = account.id
	data["name"] = account.name
	data["validationKey"] = account.validationKey
	data["withdrawalPubKey"] = hex.EncodeToString(account.withdrawalPubKey.Marshal())
	data["baseAccountPath"] = account.basePath
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
		account.id, err = uuid.Parse(val.(string))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("could not find var: id")
	}

	// name
	if val, exists := v["name"]; exists {
		account.name = val.(string)
	} else {
		return fmt.Errorf("could not find var: id")
	}

	// base path
	if val, exists := v["baseAccountPath"]; exists {
		account.basePath = val.(string)
	} else {
		return fmt.Errorf("could not find var: baseAccountPath")
	}

	// validation key
	if val, exists := v["validationKey"]; exists {
		byts, err := json.Marshal(val)
		if err != nil {
			return err
		}
		key := &core.HDKey{}
		err = json.Unmarshal(byts, key)
		if err != nil {
			return err
		}
		account.validationKey = key
	} else {
		return fmt.Errorf("could not find var: key")
	}

	// withdrawal pub Key
	if val, exists := v["withdrawalPubKey"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		account.withdrawalPubKey, err = e2types.BLSPublicKeyFromBytes(byts)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("could not find var: key")
	}

	return nil
}

func NewValidatorAccount(name string,
	validationKey *core.HDKey,
	withdrawalPubKey e2types.PublicKey,
	basePath string,
	context *core.WalletContext) (*HDAccount, error) {
	return &HDAccount{
		name:             name,
		id:               uuid.New(),
		validationKey:    validationKey,
		withdrawalPubKey: withdrawalPubKey,
		basePath:         basePath,
		context:          context,
	}, nil
}

// ID provides the ID for the account.
func (account *HDAccount) ID() uuid.UUID {
	return account.id
}

// Name provides the name for the account.
func (account *HDAccount) Name() string {
	return account.name
}

// BasePath provides the basePth of the account.
func (account *HDAccount) BasePath() string {
	return account.basePath
}

// ValidatorPublicKey provides the public key for the account.
func (account *HDAccount) ValidatorPublicKey() e2types.PublicKey {
	return account.validationKey.PublicKey()
}

// WithdrawalPublicKey provides the public key for the account.
func (account *HDAccount) WithdrawalPublicKey() e2types.PublicKey {
	return account.withdrawalPubKey
}

// Sign signs data with the account.
func (account *HDAccount) ValidationKeySign(data []byte) (e2types.Signature, error) {
	return account.validationKey.Sign(data)
}

// Get Deposit Data
func (account *HDAccount) GetDepositData() (map[string]interface{}, error) {
	depositData, root, err := eth1_deposit.DepositData(account.validationKey, account.withdrawalPubKey.Marshal(), eth1_deposit.MaxEffectiveBalanceInGwei)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"amount":                depositData.GetAmount(),
		"publicKey":             hex.EncodeToString(depositData.GetPublicKey()),
		"signature":             hex.EncodeToString(depositData.GetSignature()),
		"withdrawalCredentials": hex.EncodeToString(depositData.GetWithdrawalCredentials()),
		"depositDataRoot":       hex.EncodeToString(root[:]),
	}, nil
}

func (account *HDAccount) SetContext(ctx *core.WalletContext) {
	account.context = ctx
}
