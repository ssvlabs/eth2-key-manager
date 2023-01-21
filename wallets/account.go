package wallets

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/bloxapp/eth2-key-manager/core"
	eth1deposit "github.com/bloxapp/eth2-key-manager/eth1_deposit"
)

// HDAccount represents HD account
type HDAccount struct {
	name string
	// holds the base path from which the account was derived
	// for eip2334 should be m/12381/3600/<index>
	basePath         string
	id               uuid.UUID
	validationKey    *core.HDKey
	withdrawalPubKey []byte
	context          *core.WalletContext
}

// MarshalJSON is the custom JSON marshaler
func (account *HDAccount) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = account.id
	data["name"] = account.name
	data["validationKey"] = account.validationKey
	// Withdrawal public key is equal to public key in seedless mode
	if hex.EncodeToString(account.withdrawalPubKey) != hex.EncodeToString(account.ValidatorPublicKey()) {
		data["withdrawalPubKey"] = hex.EncodeToString(account.withdrawalPubKey)
	} else {
		// In case of seedless account we return empty string
		data["withdrawalPubKey"] = ""
	}
	data["baseAccountPath"] = account.basePath
	return json.Marshal(data)
}

// UnmarshalJSON is the custom JSON unmarshaler
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
		return errors.New("could not find var: id")
	}

	// name
	if val, exists := v["name"]; exists {
		account.name = val.(string)
	} else {
		return errors.New("could not find var: name")
	}

	// base path
	if val, exists := v["baseAccountPath"]; exists {
		account.basePath = val.(string)
	} else {
		return errors.New("could not find var: baseAccountPath")
	}

	// validation key
	if val, exists := v["validationKey"]; exists {
		byts, err := json.Marshal(val)
		if err != nil {
			return err
		}

		key := &core.HDKey{}
		if err := json.Unmarshal(byts, key); err != nil {
			return err
		}
		account.validationKey = key
	} else {
		return errors.New("could not find var: validationKey")
	}

	// withdrawal pub Key
	if val, exists := v["withdrawalPubKey"]; exists {
		byts, err := hex.DecodeString(val.(string))
		if err != nil {
			return err
		}
		account.withdrawalPubKey = byts
		if err != nil {
			return err
		}
	} else {
		return errors.New("could not find var: withdrawalPubKey")
	}

	return nil
}

// NewValidatorAccount is the constructor of HDAccount
func NewValidatorAccount(
	name string,
	validationKey *core.HDKey,
	withdrawalPubKey []byte,
	basePath string,
	context *core.WalletContext,
) *HDAccount {
	return &HDAccount{
		name:             name,
		id:               uuid.New(),
		validationKey:    validationKey,
		withdrawalPubKey: withdrawalPubKey,
		basePath:         basePath,
		context:          context,
	}
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
func (account *HDAccount) ValidatorPublicKey() []byte {
	return account.validationKey.PublicKey().Serialize()
}

// WithdrawalPublicKey provides the public key for the account.
func (account *HDAccount) WithdrawalPublicKey() []byte {
	return account.withdrawalPubKey
}

// ValidationKeySign signs data with the account.
func (account *HDAccount) ValidationKeySign(data []byte) ([]byte, error) {
	return account.validationKey.Sign(data)
}

// GetDepositData returns deposit data
func (account *HDAccount) GetDepositData() (map[string]interface{}, error) {
	depositData, root, err := eth1deposit.DepositData(
		account.validationKey,
		account.withdrawalPubKey,
		account.context.Storage.Network(),
		eth1deposit.MaxEffectiveBalanceInGwei,
	)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"amount":                 depositData.Amount,
		"publicKey":              strings.TrimPrefix(depositData.PublicKey.String(), "0x"),
		"signature":              strings.TrimPrefix(depositData.Signature.String(), "0x"),
		"withdrawalCredentials":  hex.EncodeToString(depositData.WithdrawalCredentials),
		"depositDataRoot":        hex.EncodeToString(root[:]),
		"depositContractAddress": account.context.Storage.Network().DepositContractAddress(),
	}, nil
}

// SetContext is the context setter
func (account *HDAccount) SetContext(ctx *core.WalletContext) {
	account.context = ctx
}
