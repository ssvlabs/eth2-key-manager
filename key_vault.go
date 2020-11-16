package eth2keymanager

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	e2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/wallet_hd"
)

var initBLSOnce sync.Once

// initBLS initializes BLS ONLY ONCE!
func initBLS() error {
	var err error
	var wg sync.WaitGroup
	initBLSOnce.Do(func() {
		wg.Add(1)
		err = e2types.InitBLS()
		wg.Done()
	})
	wg.Wait()
	return err
}

func InitCrypto() {
	// !!!VERY IMPORTANT!!!
	if err := initBLS(); err != nil {
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

	// update wallet context
	wallet := wallet_hd.NewHDWallet(context)

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
