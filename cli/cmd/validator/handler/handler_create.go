package handler

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/validator/flag"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
)

// ValidatorConfig represents the validator config data
type ValidatorConfig struct {
	UUID    string                 `json:"uuid"`
	Crypto  map[string]interface{} `json:"crypto"`
	PubKey  string                 `json:"pubkey"`
	Path    string                 `json:"path"`
	Version uint                   `json:"version"`
}

// Create is the handler to create validator(s).
func (h *Handler) Create(cmd *cobra.Command, args []string) error {
	// Get seeds count
	seedsCount, err := flag.GetSeedsCountFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get seeds count flag value")
	}

	// Get validators per seed number
	validatorsPerSeed, err := flag.GetValidatorsPerSeedFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get validators per seed flag value")
	}

	// Get wallet address
	walletAddress, err := flag.GetWalletAddressFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get wallet address flag value")
	}

	// Get wallet private key
	walletPrivateKey, err := flag.GetWalletPrivateKeyFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get wallet private key flag value")
	}

	// Get web3 address
	web3Addr, err := flag.GetWeb3AddrFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to get web3 address flag value")
	}

	rpcClient, err := rpc.Dial(web3Addr)
	if err != nil {
		return err
	}
	client := ethclient.NewClient(rpcClient)

	address, err := hex.DecodeString(walletAddress)
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode the given wallet address")
	}

	res, err := client.BalanceAt(context.Background(), common.BytesToAddress(address), nil)
	if err != nil {
		return err
	}
	fmt.Printf("res %#v\n", res.String())

	store := in_memory.NewInMemStore()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	if _, err = eth2keymanager.NewKeyVault(options); err != nil {
		return errors.Wrap(err, "failed to create key vault.")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	// Generate seed
	encryptor := keystorev4.New()
	var seedToAccounts []ValidatorConfig
	for i := 0; i < seedsCount; i++ {
		entropy, err := core.GenerateNewEntropy()
		if err != nil {
			return errors.Wrap(err, "failed to generate entropy")
		}

		generatedSeed, err := core.SeedFromEntropy(entropy, "")
		if err != nil {
			return errors.Wrap(err, "failed to generate seed from entropy")
		}

		// Create accounts (validators)
		for j := 0; j < validatorsPerSeed; j++ {
			account, err := wallet.CreateValidatorAccount(generatedSeed, &j)
			if err != nil {
				return errors.Wrapf(err, "failed to create validator account")
			}

			cryptoFields, err := encryptor.Encrypt(generatedSeed, "")
			if err != nil {
				return errors.Wrap(err, "could not encrypt seed phrase into keystore")
			}

			seedToAccounts = append(seedToAccounts, ValidatorConfig{
				UUID:    uuid.New().String(),
				PubKey:  hex.EncodeToString(account.ValidatorPublicKey().Marshal()),
				Path:    core.BaseEIP2334Path + account.BasePath(),
				Version: encryptor.Version(),
				Crypto:  cryptoFields,
			})
		}
	}

	h.printer.JSON(seedToAccounts)
	h.printer.JSON(walletPrivateKey)
	return nil
}
