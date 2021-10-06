package handler

import (
	"archive/zip"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	contracts "github.com/prysmaticlabs/prysm/contracts/deposit-contract"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/spf13/cobra"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/validator/flag"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/encryptor/keystorev4"
	eth1deposit "github.com/bloxapp/eth2-key-manager/eth1_deposit"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
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

	// Initialize connection with web3 API
	rpcClient, err := rpc.Dial(web3Addr)
	if err != nil {
		return errors.Wrap(err, "failed to create connection with web3 API")
	}
	defer rpcClient.Close()
	client := ethclient.NewClient(rpcClient)

	// Fetch wallet balance
	walletBalance, err := h.getWalletBalance(client, walletAddress)
	if err != nil {
		return errors.Wrap(err, "failed to get wallet balance")
	}

	store := inmemory.NewInMemStore(h.network)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	if _, err = eth2keymanager.NewKeyVault(options); err != nil {
		return errors.Wrap(err, "failed to create key vault.")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	// Check balance
	minBalance := big.NewInt(0).
		SetUint64(eth1deposit.MaxEffectiveBalanceInGwei * uint64(seedsCount) * uint64(validatorsPerSeed))
	minBalance = minBalance.Mul(minBalance, big.NewInt(1e9))
	if walletBalance.Cmp(minBalance) < 0 {
		return errors.New("insufficient funds for transfer")
	}

	// Create deposit contract client
	depositContract, err := contracts.NewDepositContract(common.HexToAddress(store.Network().DepositContractAddress()), client)
	if err != nil {
		return err
	}

	// Create transaction options
	txOpts, err := buildTransactionOpts(walletPrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to build transaction options")
	}

	// Generate seed
	encryptor := keystorev4.New()
	seedToAccounts := make(map[string][]ValidatorConfig)
	defer func() {
		if err := h.writeResultToFiles(seedToAccounts); err != nil {
			h.printer.Error(err)
			if err := h.printer.JSON(seedToAccounts); err != nil {
				h.printer.Error(err)
			}
		}
	}()
	for i := 0; i < seedsCount; i++ {
		entropy, err := core.GenerateNewEntropy()
		if err != nil {
			return errors.Wrap(err, "failed to generate entropy")
		}

		mnemonic, err := core.EntropyToMnemonic(entropy)
		if err != nil {
			return errors.Wrap(err, "failed to generate mnemonic from entropy")
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

			// Make transaction
			if err := h.makeTransaction(depositContract, txOpts, account); err != nil {
				return errors.Wrap(err, "failed to make deposit")
			}

			if _, ok := seedToAccounts[mnemonic]; !ok {
				seedToAccounts[mnemonic] = []ValidatorConfig{}
			}
			seedToAccounts[mnemonic] = append(seedToAccounts[mnemonic], ValidatorConfig{
				UUID:    uuid.New().String(),
				PubKey:  hex.EncodeToString(account.ValidatorPublicKey()),
				Path:    store.Network().FullPath(account.BasePath()),
				Version: encryptor.Version(),
				Crypto:  cryptoFields,
			})
		}
	}

	return nil
}

func (h *Handler) getWalletBalance(client *ethclient.Client, walletAddr string) (*big.Int, error) {
	address, err := hex.DecodeString(walletAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode the given wallet address")
	}

	res, err := client.BalanceAt(context.Background(), common.BytesToAddress(address), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get balance of the wallet")
	}

	return res, nil
}

func (h *Handler) makeTransaction(depositContract *contracts.DepositContract, txOpts *bind.TransactOpts, account core.ValidatorAccount) error {
	// Get deposit data for account
	depositData, err := account.GetDepositData()
	if err != nil {
		return errors.Wrap(err, "failed to get deposit data")
	}

	// Prepare public key
	publicKey, err := hex.DecodeString(depositData["publicKey"].(string))
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode account public key")
	}

	// Prepare withdrawal credentials
	withdrawalCredentials, err := hex.DecodeString(depositData["withdrawalCredentials"].(string))
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode account withdrawal credentials")
	}

	// Prepare signature
	signature, err := hex.DecodeString(depositData["signature"].(string))
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode account signature")
	}

	// Prepare data root
	depositDataRoot, err := hex.DecodeString(depositData["depositDataRoot"].(string))
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode account deposit data root")
	}

	// Sent deposit contract
	tx, err := depositContract.Deposit(
		txOpts,
		publicKey,
		withdrawalCredentials,
		signature,
		bytesutil.ToBytes32(depositDataRoot),
	)
	if err != nil {
		return errors.Wrap(err, "unable to send transaction to contract")
	}

	h.printer.Text(fmt.Sprintf("Monitor your transaction on Etherscan here https://goerli.etherscan.io/tx/0x%x", tx.Hash()))
	return nil
}

func (h *Handler) writeResultToFiles(results map[string][]ValidatorConfig) error {
	if len(results) == 0 {
		return errors.New("no results to store")
	}

	// Create zip file
	fileName := fmt.Sprintf("validators_%d_%d", time.Now().Unix(), len(results))
	out, cleanup, err := h.resultWriterFactory(fileName)
	if err != nil {
		return errors.Wrap(err, "failed to create result writer")
	}
	defer cleanup()

	// Create a new zip archive.
	w := zip.NewWriter(out)

	// Put results into archive
	for mnemonic, validators := range results {
		for _, validator := range validators {
			// Create file
			f, err := w.Create(mnemonic + "/keystore-" + strings.ReplaceAll(validator.Path, "/", "_") + ".json")
			if err != nil {
				return errors.Wrap(err, "failed to create result file")
			}

			// Put result into file
			if err := json.NewEncoder(f).Encode(validator); err != nil {
				return errors.Wrap(err, "failed to write result into file")
			}
		}
	}

	// Close archive
	if err := w.Close(); err != nil {
		h.printer.Error(err)
	}

	return nil
}

func buildTransactionOpts(privateKey string) (*bind.TransactOpts, error) {
	// User inputs private key, sign tx with private key
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode private key")
	}

	txOps, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	if err != nil {
		return nil, err
	}
	txOps.Value = new(big.Int).Mul(big.NewInt(int64(eth1deposit.MaxEffectiveBalanceInGwei)), big.NewInt(1e9))
	txOps.GasLimit = 500000
	txOps.Context = context.Background()
	return txOps, nil
}
