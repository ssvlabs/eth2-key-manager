package KeyVault

import (
	"crypto/rand"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/encryptors"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type WalletOptions struct {
	encryptor wtypes.Encryptor
	password []byte
	storage interface{}
	enableSimpleSigner bool
	seed []byte
	portfolioLockPolicy core.LockablePolicy
}

func (options *WalletOptions)SetPortfolioLockPolicy(lockPolicy core.LockablePolicy) *WalletOptions {
	options.portfolioLockPolicy = lockPolicy
	return options
}

func (options *WalletOptions)SetEncryptor(encryptor wtypes.Encryptor) *WalletOptions {
	options.encryptor = encryptor
	return options
}

func (options *WalletOptions)SetStorage(storage interface{}) *WalletOptions {
	options.storage = storage
	return options
}

func (options *WalletOptions)SetPassword(password string) *WalletOptions {
	options.password = []byte(password)
	return options
}

func (options *WalletOptions)EnableSimpleSigner(val bool) *WalletOptions {
	options.enableSimpleSigner = true
	return options
}

func (options *WalletOptions)SetSeed(seed []byte) *WalletOptions {
	options.seed = seed
	return options
}

func (options *WalletOptions) GenerateSeed() error {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)

	options.SetSeed(seed)

	return err
}

func (options *WalletOptions) setNoEncryptor() *WalletOptions{
	return options.SetEncryptor(encryptors.NewPlainTextEncryptor())
}

func (options *WalletOptions) setNoLockPolicy() *WalletOptions{
	return options.SetPortfolioLockPolicy(&core.NoLockPolicy{})
}