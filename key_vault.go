package KeyVault

import (
	"crypto/rand"
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
	id          uuid.UUID
	//indexMapper map[string]uuid.UUID
	Context     *core.WalletContext
	//key         *core.MasterDerivableKey
	walletId 	uuid.UUID
}

type NotExistError struct {
	desc string
}

func (e *NotExistError) Error() string {
	return fmt.Sprintf("%s", e.desc)
}

func OpenKeyVault(options *WalletOptions) (*KeyVault, error) {
	// storage
	storage, err := setupStorage(options)
	if err != nil {
		return nil, err
	}

	wallet, err := storage.OpenWallet()
	if err != nil {
		return nil, err
	}

	return completeVaultSetup(options, wallet)
}

func ImportKeyVault(options *WalletOptions) (*KeyVault, error) {
	// storage
	storage, err := setupStorage(options)
	if err != nil {
		return nil, err
	}

	// key
	if options.seed == nil {
		return nil, fmt.Errorf("no seed was provided")
	}
	err = storage.SecurelySavePortfolioSeed(options.seed)
	if err != nil {
		return nil, err
	}
	key, err := core.MasterKeyFromSeed(storage)
	if err != nil {
		return nil, err
	}

	return completeVaultSetup(options, wallet_hd.NewHDWallet(key, nil))
}

func NewKeyVault(options *WalletOptions) (*KeyVault, error) {
	// storage
	storage, err := setupStorage(options)
	if err != nil {
		return nil, err
	}

	// key
	seed, err := storage.SecurelyFetchPortfolioSeed()
	if err != nil || len(seed) == 0 {
		seed, err = saveNewSeed(storage)
	}
	if err != nil {
		return nil, err
	}

	key, err := core.MasterKeyFromSeed(storage)
	if err != nil {
		return nil, err
	}

	return completeVaultSetup(options, wallet_hd.NewHDWallet(key, nil))
}

func (kv *KeyVault) Wallet() (core.Wallet, error) {
	return kv.Context.Storage.OpenWallet()
}

func completeVaultSetup(options *WalletOptions, wallet core.Wallet) (*KeyVault, error) {
	// wallet Context
	context := &core.WalletContext{
		Storage: options.storage.(core.Storage),
	}

	// update wallet context
	wallet.SetContext(context)

	ret := &KeyVault{
		id:          uuid.New(),
		Context:     context,
		walletId:    wallet.ID(),
	}

	err := options.storage.(core.Storage).SaveWallet(wallet)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func setupStorage(options *WalletOptions) (core.Storage, error) {
	if _, ok := options.storage.(core.Storage); !ok {
		return nil, fmt.Errorf("storage does not implement core.Storage")
	} else {
		if options.encryptor != nil && options.password != nil {
			options.storage.(core.Storage).SetEncryptor(options.encryptor, options.password)
		}
	}

	return options.storage.(core.Storage), nil
}

func saveNewSeed(storage core.Storage) ([]byte, error) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, err
	}
	err = storage.SecurelySavePortfolioSeed(seed)
	if err != nil {
		return nil, err
	}

	return seed, nil
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
