package wallet_hd

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// according to https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
const (
	BaseAccountPath   = "/%d"
	WithdrawalKeyPath = BaseAccountPath + "/0"
	ValidatorKeyPath  = WithdrawalKeyPath + "/0"
)

// an hierarchical deterministic wallet
type HDWallet struct {
	id          uuid.UUID
	walletType  core.WalletType
	indexMapper map[string]uuid.UUID
	context     *core.WalletContext
}

func NewHDWallet(context *core.WalletContext) *HDWallet {
	return &HDWallet{
		id:          uuid.New(),
		walletType:  core.HDWallet,
		indexMapper: make(map[string]uuid.UUID),
		context:     context,
	}
}

// ID provides the ID for the wallet.
func (wallet *HDWallet) ID() uuid.UUID {
	return wallet.id
}

// Type provides the type of the wallet.
func (wallet *HDWallet) Type() core.WalletType {
	return wallet.walletType
}

// CreatePrivateKey creates a private key
func CreatePrivateKey(seed []byte, path string, index int) (*core.HDKey, error) {
	// create the master key
	masterKey, err := core.MasterKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}

	keyPath := fmt.Sprintf(path, index)
	key, err := masterKey.Derive(keyPath)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
// This will error if an account with the name already exists.
func (wallet *HDWallet) CreateValidatorAccount(seed []byte, name string) (core.ValidatorAccount, error) {
	if len(name) == 0 {
		name = fmt.Sprintf("account-%d", len(wallet.indexMapper))
	}

	validatorKey, err := CreatePrivateKey(seed, ValidatorKeyPath, len(wallet.indexMapper))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create validator key")
	}
	withdrawalKey, err := CreatePrivateKey(seed, WithdrawalKeyPath, len(wallet.indexMapper))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create withdrawal key")
	}
	baseAccountPath := fmt.Sprintf(BaseAccountPath, len(wallet.indexMapper))

	// create ret account
	ret, err := NewValidatorAccount(
		name,
		validatorKey,
		withdrawalKey.PublicKey(),
		baseAccountPath,
		wallet.context,
	)
	if err != nil {
		return nil, err
	}

	validatorPublicKey := hex.EncodeToString(ret.ValidatorPublicKey().Marshal())
	// register new wallet and save portfolio
	reset := func() {
		delete(wallet.indexMapper, validatorPublicKey)
	}
	wallet.indexMapper[validatorPublicKey] = ret.ID()
	err = wallet.context.Storage.SaveAccount(ret)
	if err != nil {
		reset()
		return nil, err
	}
	err = wallet.context.Storage.SaveWallet(wallet)
	if err != nil {
		reset()
		return nil, err
	}

	return ret, nil
}

// Accounts provides all accounts in the wallet.
func (wallet *HDWallet) Accounts() <-chan core.ValidatorAccount {
	ch := make(chan core.ValidatorAccount, 1024) // TODO - handle more? change from chan?
	go func() {
		for pubKey := range wallet.indexMapper {
			id := wallet.indexMapper[pubKey]
			account, err := wallet.AccountByID(id)
			if err != nil {
				continue
			}
			ch <- account
		}
		close(ch)
	}()

	return ch
}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByID(id uuid.UUID) (core.ValidatorAccount, error) {
	ret, err := wallet.context.Storage.OpenAccount(id)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	ret.SetContext(wallet.context)
	return ret, nil
}

func (wallet *HDWallet) SetContext(ctx *core.WalletContext) {
	wallet.context = ctx
}

// AccountByPublicKey provides a single account from the wallet given its public key.
// This will error if the account is not found.
func (wallet *HDWallet) AccountByPublicKey(pubKey string) (core.ValidatorAccount, error) {
	id, exists := wallet.indexMapper[pubKey]
	if !exists {
		return nil, fmt.Errorf("account not found")
	}
	return wallet.AccountByID(id)
}
