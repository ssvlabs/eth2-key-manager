package KeyVault

import (
	"fmt"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"log"
	"sync"

	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
	e2types "github.com/wealdtech/go-eth2-types/v2"
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

func init() {
	// !!!VERY IMPORTANT!!!
	if err := initBLS(); err != nil {
		log.Fatal(err)
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

	err = options.storage.(core.Storage).SaveWallet(wallet)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func setupStorage(options *KeyVaultOptions) (core.Storage, error) {
	if _, ok := options.storage.(core.Storage); !ok {
		return nil, fmt.Errorf("storage does not implement core.Storage")
	} else {
		if options.encryptor != nil && options.password != nil {
			options.storage.(core.Storage).SetEncryptor(options.encryptor, options.password)
		}
	}

	return options.storage.(core.Storage), nil
}

func GenerateNewSeed() ([]byte, error) {
	seed, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}

	return seed, nil
}

func SeedToMnemonic(seed []byte) (string, error) {
	mnemonic, err := bip39.NewMnemonic(seed)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func SeedFromMnemonic(mnemonic string) ([]byte, error) {
	seed, err := bip39.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return seed, nil
}
