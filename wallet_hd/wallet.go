package wallet_hd

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/bloxapp/eth2-key-manager/core"
)

// according to https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
const (
	BaseAccountPath   = "/%d"
	WithdrawalKeyPath = BaseAccountPath + "/0"
	ValidatorKeyPath  = WithdrawalKeyPath + "/0"
)

// Predefined errors
var (
	// ErrAccountNotFound is the error when account not found
	ErrAccountNotFound = errors.New("account not found")
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

// GetNextAccountIndex provides next index to create account at.
func (wallet *HDWallet) GetNextAccountIndex() int {
	if len(wallet.indexMapper) == 0 {
		return 0
	}
	accounts := wallet.Accounts()
	index, _ := strconv.ParseInt(accounts[0].BasePath()[1:], 0, 64)
	return int(index) + 1
}

// CreateValidatorKey creates a new validation (validator) key pair in the wallet.
func (wallet *HDWallet) CreateValidatorAccount(seed []byte, indexPointer *int) (core.ValidatorAccount, error) {
	// Resolve index to create account at
	var index int
	if indexPointer != nil {
		index = *indexPointer
	} else {
		index = wallet.GetNextAccountIndex()
	}
	name := fmt.Sprintf("account-%d", index)

	// Create the master key based on the seed and network.
	key, err := core.MasterKeyFromSeed(seed, wallet.context.Storage.Network())
	if err != nil {
		return nil, err
	}

	baseAccountPath := fmt.Sprintf(BaseAccountPath, index)

	// Create validator key
	validatorPath := fmt.Sprintf(ValidatorKeyPath, index)
	validatorKey, err := key.Derive(validatorPath)
	if err != nil {
		return nil, err
	}

	// Create withdrawal key
	withdrawalPath := fmt.Sprintf(WithdrawalKeyPath, index)
	withdrawalKey, err := key.Derive(withdrawalPath)
	if err != nil {
		return nil, err
	}

	// Create ret account
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

	// Register new wallet and save portfolio
	reset := func() {
		delete(wallet.indexMapper, validatorPublicKey)
	}
	wallet.indexMapper[validatorPublicKey] = ret.ID()

	// Store account
	if err = wallet.context.Storage.SaveAccount(ret); err != nil {
		reset()
		return nil, err
	}

	// Store wallet
	err = wallet.context.Storage.SaveWallet(wallet)
	if err != nil {
		reset()
		return nil, err
	}

	return ret, nil
}

func (wallet *HDWallet) DeleteAccountByPublicKey(pubKey string) error {
	account, err := wallet.AccountByPublicKey(pubKey)
	if err != nil {
		return errors.Wrap(err, "failed to get account by public key")
	}

	err = wallet.context.Storage.DeleteAccount(account.ID())
	if err != nil {
		return errors.Wrap(err, "failed to delete account from store")
	}
	delete(wallet.indexMapper, pubKey)
	err = wallet.context.Storage.SaveWallet(wallet)
	if err != nil {
		return errors.Wrap(err, "failed to save wallet")
	}
	return nil
}

// Accounts provides all accounts in the wallet.
func (wallet *HDWallet) Accounts() []core.ValidatorAccount {
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
		return nil, ErrAccountNotFound
	}
	return wallet.AccountByID(id)
}
