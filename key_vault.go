package eth2keymanager

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/wallets/hd"
	"github.com/bloxapp/eth2-key-manager/wallets/nd"

	// Force import a transitive dependency to fix an ambiguous import error.
	// See https://github.com/btcsuite/btcd/issues/1839
	_ "github.com/btcsuite/btcd/btcec/v2"
)

// InitCrypto initializes cryptography
func InitCrypto() {
	// !!!VERY IMPORTANT!!!
	if err := core.InitBLS(); err != nil {
		logrus.Fatal(err)
	}
}

// KeyVault is an EIP 2333,2334,2335 compliant hierarchical deterministic portfolio
// https://eips.ethereum.org/EIPS/eip-2333
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
// https://eips.ethereum.org/EIPS/eip-2335
type KeyVault struct {
	Context  *core.WalletContext
	walletID uuid.UUID
}

// Wallet returns wallet
func (kv *KeyVault) Wallet() (core.Wallet, error) {
	return kv.Context.Storage.OpenWallet()
}

// OpenKeyVault opens an existing KeyVault (and wallet) from memory
func OpenKeyVault(options *KeyVaultOptions) (*KeyVault, error) {
	storage, err := setupStorage(options)
	if err != nil {
		return nil, err
	}

	// wallet Context
	context := &core.WalletContext{
		Storage: storage,
	}

	// try and open a wallet
	wallet, err := storage.OpenWallet()
	if err != nil {
		return nil, err
	}

	return &KeyVault{
		Context:  context,
		walletID: wallet.ID(),
	}, nil
}

// NewKeyVault creates a new wallet (with new ids) and will save it to storage
// Import and New are the same action.
func NewKeyVault(options *KeyVaultOptions) (*KeyVault, error) {
	storage, err := setupStorage(options)
	if err != nil {
		return nil, err
	}

	// wallet Context
	context := &core.WalletContext{
		Storage: storage,
	}

	// create wallet
	var wallet core.Wallet
	if options.walletType == core.NDWallet {
		wallet = nd.NewWallet(context)
	} else { // ND wallet by default
		wallet = hd.NewWallet(context)
	}

	ret := &KeyVault{
		Context:  context,
		walletID: wallet.ID(),
	}

	storage, ok := options.storage.(core.Storage)
	if !ok {
		return nil, errors.Errorf("unexpected storage type %T", options.storage)
	}

	if err := storage.SaveWallet(wallet); err != nil {
		return nil, err
	}

	return ret, nil
}

func setupStorage(options *KeyVaultOptions) (core.Storage, error) {
	storage, ok := options.storage.(core.Storage)
	if !ok {
		return nil, errors.New("storage does not implement core.Storage")
	}

	if options.encryptor != nil && options.password != nil {
		storage.SetEncryptor(options.encryptor, options.password)
	}

	return storage, nil
}

// This function calls anyway when this package is imported by someone.
// No needed to call this in each function.
func init() {
	InitCrypto()
}
