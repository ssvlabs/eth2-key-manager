package KeyVault

import (
	"crypto/rand"
	"github.com/bloxapp/KeyVault/encryptors"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type PortfolioOptions struct {
	encryptor wtypes.Encryptor
	password []byte
	storage interface{}
	enableSimpleSigner bool
	seed []byte
}

func (options *PortfolioOptions)SetEncryptor(encryptor wtypes.Encryptor) *PortfolioOptions {
	options.encryptor = encryptor
	return options
}

func (options *PortfolioOptions)SetStorage(storage interface{}) *PortfolioOptions {
	options.storage = storage
	return options
}

func (options *PortfolioOptions)SetPassword(password string) *PortfolioOptions {
	options.password = []byte(password)
	return options
}

func (options *PortfolioOptions)EnableSimpleSigner(val bool) *PortfolioOptions {
	options.enableSimpleSigner = true
	return options
}

func (options *PortfolioOptions)SetSeed(seed []byte) *PortfolioOptions {
	options.seed = seed
	return options
}

func (options *PortfolioOptions) GenerateSeed() error {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)

	options.SetSeed(seed)

	return err
}

func (options *PortfolioOptions) setNoEncryptor() *PortfolioOptions {
	return options.SetEncryptor(encryptors.NewPlainTextEncryptor())
}
