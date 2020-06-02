package hd

import (
	"fmt"
	core "github.com/bloxapp/KeyVault/core"
	"github.com/google/uuid"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"github.com/wealdtech/go-indexer"
)

// This is an EIP 2333,2334,2335 compliant hierarchical deterministic wallet
//https://eips.ethereum.org/EIPS/eip-2333
//https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2334.md
//https://eips.ethereum.org/EIPS/eip-2335
type HDPortfolio struct {
	storage core.PortfolioStorage
	encryptor types.Encryptor
	seed *core.EncryptableSeed
	lockPolicy core.LockablePolicy
	walletsIndexer indexer.Index // maps indexs <> names
	walletIds []uuid.UUID
	lockPassword []byte
}

// CreateAccount creates a new account in the wallet.
// This will error if an account with the name already exists.
// Will push to the new wallet the lock policy
func (portfolio *HDPortfolio) CreateWallet(name string) (core.Wallet, error) {

}

// Accounts provides all accounts in the wallet.
func (portfolio *HDPortfolio) Wallets() (<-chan core.Wallet,error) {
	if portfolio.IsLocked() {
		return nil,fmt.Errorf("portfolio is locked")
	}

	ch := make (chan core.Wallet,1024) // TODO - handle more?
	go func() {
		for i := range portfolio.walletIds {
			id := portfolio.walletIds[i]
			wallet,err := portfolio.WalletByID(id)
			if err != nil {
				continue
			}
			ch <- wallet
		}

		if portfolio.lockPolicy.LockAfterOperation(core.FetchData) {
			portfolio.Lock()
		}
	}()

	return ch,nil
}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (portfolio *HDPortfolio) WalletByID(id uuid.UUID) (core.Wallet, error) {
	if portfolio.IsLocked() {
		return nil,fmt.Errorf("portfolio is locked")
	}
	defer func() {
		if portfolio.lockPolicy.LockAfterOperation(core.FetchData) {
			portfolio.Lock()
		}
	}()

}

// AccountByName provides a single account from the wallet given its name.
// This will error if the account is not found.
func (portfolio *HDPortfolio) WalletByName(name string) (core.Wallet, error) {
	if portfolio.IsLocked() {
		return nil,fmt.Errorf("portfolio is locked")
	}
	defer func() {
		if portfolio.lockPolicy.LockAfterOperation(core.FetchData) {
			portfolio.Lock()
		}
	}()
}

func (portfolio *HDPortfolio) Lock() error {
	return portfolio.seed.Encrypt(portfolio.lockPassword)
}

func (portfolio *HDPortfolio) IsLocked() bool {
	return portfolio.seed.IsEncrypted()
}

func (portfolio *HDPortfolio) Unlock(password []byte) error {
	return portfolio.seed.Decrypt(password)
}