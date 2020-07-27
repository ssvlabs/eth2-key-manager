package KeyVault

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	util "github.com/wealdtech/go-eth2-util"
	"log"
	"sync"

	"github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

const (
	BaseEIP2334Path = "m/12381/3600"
	//TODO change to /0/%d when remove portfolio from the path
	ValidatorKeyPath = "/0/0/%d"
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
	indexMapper map[string]uuid.UUID
	Context     *core.PortfolioContext
	key         *core.DerivableKey
}

type NotExistError struct {
	desc string
}

func (e *NotExistError) Error() string {
	return fmt.Sprintf("%s", e.desc)
}

func OpenKeyVault(options *PortfolioOptions) (*KeyVault, error) {
	// storage
	storage, err := setupStorage(options)
	if err != nil {
		return nil, err
	}

	bytes, err := storage.OpenPortfolioRaw()
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return nil, &NotExistError{"key vault not found"}
	}

	// portfolio Context
	context := &core.PortfolioContext{
		Storage: options.storage.(core.Storage),
	}

	ret := &KeyVault{Context: context}
	err = json.Unmarshal(bytes, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func ImportKeyVault(options *PortfolioOptions) (*KeyVault, error) {
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
	key, err := core.BaseKeyFromSeed(options.seed, storage)
	if err != nil {
		return nil, err
	}

	return completeVaultSetup(options, key)
}

func NewKeyVault(options *PortfolioOptions) (*KeyVault, error) {
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

	key, err := core.BaseKeyFromSeed(seed, storage)
	if err != nil {
		return nil, err
	}

	return completeVaultSetup(options, key)
}

func completeVaultSetup(options *PortfolioOptions, key *core.DerivableKey) (*KeyVault, error) {
	// portfolio Context
	context := &core.PortfolioContext{
		Storage: options.storage.(core.Storage),
	}

	ret := &KeyVault{
		id:          uuid.New(),
		indexMapper: make(map[string]uuid.UUID),
		Context:     context,
		key:         key,
	}

	// update Context with portfolio id
	context.PortfolioId = ret.ID()

	return ret, nil
}

func setupStorage(options *PortfolioOptions) (core.Storage, error) {
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

func CreateAccount(seed []byte, index int) ([]byte, error) {
	if seed == nil {
		return nil, fmt.Errorf("no seed was provided")
	}
	relativePath := fmt.Sprintf(ValidatorKeyPath, index)
	// TODO Validate relative path
	path := BaseEIP2334Path + relativePath
	key, err := util.PrivateKeyFromSeedAndPath(seed, path)
	if err != nil {
		return nil, err
	}

	return key.Marshal(), nil
}