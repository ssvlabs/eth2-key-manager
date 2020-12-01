package eth2keymanager

import (
	"github.com/bloxapp/eth2-key-manager/wallets/hd"

	"github.com/bloxapp/eth2-key-manager/wallets/nd"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bloxapp/eth2-key-manager/core"
)

func InitCrypto() {
	// !!!VERY IMPORTANT!!!
	if err := core.InitBLS(); err != nil {
		logrus.Fatal(err)
	}
}

// This is an EIP 2333,2334,2335 compliant hierarchical deterministic portfolio
//https://eips.ethereum.org/EIPS/eip-2333
//https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
//https://eips.ethereum.org/EIPS/eip-2335
type KeyVault struct {
	Context  *core.WalletContext
	walletId uuid.UUID
}

func (kv *KeyVault) Wallet() (core.Wallet, error) {
	return kv.Context.Storage.OpenWallet()
}

// wil try and open an existing KeyVault (and wallet) from memory
func OpenKeyVault(options *KeyVaultOptions) (*KeyVault, error) {
	InitCrypto()

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
		walletId: wallet.ID(),
	}, nil
}

// New KeyVault will create a new wallet (with new ids) and will save it to storage
// Import and New are the same action.
func NewKeyVault(options *KeyVaultOptions) (*KeyVault, error) {
	InitCrypto()

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
		wallet = nd.NewNDWallet(context)
	} else { // ND wallet by default
		wallet = hd.NewHDWallet(context)
	}

	ret := &KeyVault{
		Context:  context,
		walletId: wallet.ID(),
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
