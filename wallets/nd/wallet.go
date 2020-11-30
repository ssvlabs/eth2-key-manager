package nd

import (
	"encoding/hex"
	"sort"
	"strconv"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Predefined errors
var (
	// ErrAccountNotFound is the error when account not found
	ErrAccountNotFound = errors.New("account not found")
)

// an hierarchical deterministic wallet
type NDWallet struct {
	id          uuid.UUID
	walletType  core.WalletType
	indexMapper map[string]uuid.UUID
	context     *core.WalletContext
}

func NewNDWallet(context *core.WalletContext) *NDWallet {
	return &NDWallet{
		id:          uuid.New(),
		walletType:  core.NDWallet,
		indexMapper: make(map[string]uuid.UUID),
		context:     context,
	}
}

// ID provides the ID for the wallet.
func (wallet *NDWallet) ID() uuid.UUID {
	return wallet.id
}

// Type provides the type of the wallet.
func (wallet *NDWallet) Type() core.WalletType {
	return wallet.walletType
}

// GetNextAccountIndex provides next index to create account at.
func (wallet *NDWallet) GetNextAccountIndex() int {
	if len(wallet.indexMapper) == 0 {
		return 0
	}
	accounts := wallet.Accounts()
	index, _ := strconv.ParseInt(accounts[0].BasePath()[1:], 0, 64)
	return int(index) + 1
}

// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
func (wallet *NDWallet) CreateValidatorAccount(seed []byte, indexPointer *int) (core.ValidatorAccount, error) {
	return nil, errors.Errorf("non deterministic wallet can't create validator, please use AddValidatorAccount")
}

func (wallet *NDWallet) AddValidatorAccount(account core.ValidatorAccount) error {
	validatorPublicKey := hex.EncodeToString(account.ValidatorPublicKey().Marshal())
	wallet.indexMapper[validatorPublicKey] = account.ID()

	// Store account
	if err := wallet.context.Storage.SaveAccount(account); err != nil {
		return err
	}

	// Store wallet
	err := wallet.context.Storage.SaveWallet(wallet)
	if err != nil {
		return err
	}

	return nil
}

func (wallet *NDWallet) DeleteAccountByPublicKey(pubKey string) error {
	account, err := wallet.AccountByPublicKey(pubKey)
	if err != nil {
		return errors.Wrap(err, "failed to get account by public key")
	}

	if err := wallet.context.Storage.DeleteAccount(account.ID()); err != nil {
		return errors.Wrap(err, "failed to delete account from store")
	}
	delete(wallet.indexMapper, pubKey)

	if err := wallet.context.Storage.SaveWallet(wallet); err != nil {
		return errors.Wrap(err, "failed to save wallet")
	}
	return nil
}

// Accounts provides all accounts in the wallet.
func (wallet *NDWallet) Accounts() []core.ValidatorAccount {
	accounts := make([]core.ValidatorAccount, 0)
	for pubKey := range wallet.indexMapper {
		id := wallet.indexMapper[pubKey]
		account, err := wallet.AccountByID(id)
		if err != nil {
			continue
		}
		accounts = append(accounts, account)
	}
	sort.Slice(accounts, func(i, j int) bool {
		a, _ := strconv.ParseInt(accounts[i].BasePath()[1:], 0, 64)
		b, _ := strconv.ParseInt(accounts[j].BasePath()[1:], 0, 64)
		return a > b
	})
	return accounts
}

// AccountByID provides a nd account from the wallet given its ID.
// This will error if the account is not found.
func (wallet *NDWallet) AccountByID(id uuid.UUID) (core.ValidatorAccount, error) {
	ret, err := wallet.context.Storage.OpenAccount(id)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, ErrAccountNotFound
	}

	ret.SetContext(wallet.context)
	return ret, nil
}

func (wallet *NDWallet) SetContext(ctx *core.WalletContext) {
	wallet.context = ctx
}

// AccountByPublicKey provides a nd account from the wallet given its public key.
// This will error if the account is not found.
func (wallet *NDWallet) AccountByPublicKey(pubKey string) (core.ValidatorAccount, error) {
	id, exists := wallet.indexMapper[pubKey]
	if !exists {
		return nil, ErrAccountNotFound
	}
	return wallet.AccountByID(id)
}
