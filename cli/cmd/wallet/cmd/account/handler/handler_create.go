package handler

import (
	"encoding/hex"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	eth2keymanager "github.com/ssvlabs/eth2-key-manager"
	rootcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/ssvlabs/eth2-key-manager/core"
	"github.com/ssvlabs/eth2-key-manager/stores/inmemory"
)

// CreateAccountFlagValues keeps all collected values for seed and seedless modes
type CreateAccountFlagValues struct {
	index            int
	seed             string
	seedBytes        []byte
	privateKeys      [][]byte
	accumulate       bool
	responseType     rootcmd.ResponseType
	highestSources   []uint64
	highestTargets   []uint64
	highestProposals []uint64
	network          core.Network
}

// Create creates a new wallet account(s) and prints the storage.
func (h *Account) Create(cmd *cobra.Command, args []string) error {
	err := core.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	accountFlags, err := CollectAccountFlags(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to collect account flags")
	}

	err = h.BuildAccounts(accountFlags)
	if err != nil {
		return errors.Wrap(err, "failed to build accounts")
	}
	return nil
}

// BuildAccounts builds accounts based on account flags.
func (h *Account) BuildAccounts(accountFlags *CreateAccountFlagValues) error {
	// Initialize store
	store := inmemory.NewInMemStore(accountFlags.network)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	// Create new key vault
	_, err := eth2keymanager.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	if accountFlags.accumulate {
		for i := 0; i <= accountFlags.index; i++ {
			err := GenerateAccounts(wallet, store, i, accountFlags)
			if err != nil {
				return err
			}
		}
	} else {
		err := GenerateAccounts(wallet, store, accountFlags.index, accountFlags)
		if err != nil {
			return err
		}
	}

	if accountFlags.responseType == rootcmd.StorageResponseType {
		// marshal storage
		bytes, err := store.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "failed to JSON marshal storage")
		}

		h.printer.Text(hex.EncodeToString(bytes))
		return nil
	}

	var accounts []map[string]string
	for _, a := range wallet.Accounts() {
		var withdrawalPubKey string
		if len(accountFlags.privateKeys) > 0 {
			withdrawalPubKey = ""
		} else {
			withdrawalPubKey = hex.EncodeToString(a.WithdrawalPublicKey())
		}
		accObj := map[string]string{
			"id":               a.ID().String(),
			"name":             a.Name(),
			"validationPubKey": hex.EncodeToString(a.ValidatorPublicKey()),
			"withdrawalPubKey": withdrawalPubKey,
		}
		accounts = append(accounts, accObj)
	}

	if accountFlags.accumulate || len(accountFlags.privateKeys) > 1 {
		err = h.printer.JSON(accounts)
	} else if len(accounts) > 0 {
		err = h.printer.JSON(accounts[0])
	}
	if err != nil {
		return errors.Wrap(err, "failed to print accounts JSON")
	}
	return nil
}

// CollectAccountFlags returns collected flags for seed and seedless modes
func CollectAccountFlags(cmd *cobra.Command) (*CreateAccountFlagValues, error) {
	accountFlagValues := CreateAccountFlagValues{}

	// Seedless mode
	if cmd.Flags().Changed(flag.GetPrivateKeyFlagName()) {
		// Get privateKey flag value.
		privateKeyValues, err := flag.GetPrivateKeyFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the private key flag value")
		}

		privateKeys := strings.Split(privateKeyValues, ",")

		for _, pk := range privateKeys {
			privateKeyBytes, err := hex.DecodeString(pk)
			if err != nil {
				return nil, errors.Wrap(err, "failed to HEX decode private-key")
			}
			accountFlagValues.privateKeys = append(accountFlagValues.privateKeys, privateKeyBytes)
		}
	} else {
		// Get seed flag value.
		seedFlagValue, err := rootcmd.GetSeedFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the seed flag value")
		}
		accountFlagValues.seed = seedFlagValue

		// Get seed bytes
		seedBytes, err := hex.DecodeString(seedFlagValue)
		if err != nil {
			return nil, errors.Wrap(err, "failed to HEX decode seed")
		}
		accountFlagValues.seedBytes = seedBytes

		// Get accumulate flag value.
		accumulateFlagValue, err := rootcmd.GetAccumulateFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the accumulate flag value")
		}
		accountFlagValues.accumulate = accumulateFlagValue
	}

	// Get index flag value.
	indexFlagValue, err := rootcmd.GetIndexFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the index flag value")
	}
	accountFlagValues.index = indexFlagValue

	// Get response-type flag value.
	responseType, err := rootcmd.GetResponseTypeFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the response type value")
	}
	accountFlagValues.responseType = responseType

	// Get HighestSource flag value.
	highestSources, err := flag.GetHighestSourceFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}
	accountFlagValues.highestSources = highestSources

	// Get HighestTarget flag value.
	highestTargets, err := flag.GetHighestTargetFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}
	accountFlagValues.highestTargets = highestTargets

	// Get HighestProposal flag value.
	highestProposals, err := flag.GetHighestProposalFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}
	accountFlagValues.highestProposals = highestProposals

	// Validate highest attestation/proposal values
	highestValuesError := ValidateHighestValues(accountFlagValues)
	if highestValuesError != nil {
		return nil, highestValuesError
	}

	// Get network flag value.
	network, err := rootcmd.GetNetworkFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the network flag value")
	}
	accountFlagValues.network = network

	return &accountFlagValues, nil
}

// GenerateAccounts generates account by index using provided account flags
func GenerateAccounts(wallet core.Wallet, store *inmemory.InMemStore, index int, accountFlags *CreateAccountFlagValues) error {
	var acc core.ValidatorAccount
	var err error

	if len(accountFlags.privateKeys) > 0 {
		for i, pk := range accountFlags.privateKeys {
			indexToCreateAccountAt := accountFlags.index + i
			acc, err = wallet.CreateValidatorAccountFromPrivateKey(pk, &indexToCreateAccountAt)
			if err != nil {
				return errors.Wrap(err, "failed to create validator account from private key")
			}

			err = SaveHighestData(acc, store, accountFlags, i)
			if err != nil {
				return errors.Wrap(err, "failed to save highest att/proposal data for account")
			}
		}
	} else {
		acc, err = wallet.CreateValidatorAccount(accountFlags.seedBytes, &index)
		if err != nil {
			return errors.Wrap(err, "failed to create validator account")
		}

		err = SaveHighestData(acc, store, accountFlags, index)
		if err != nil {
			return errors.Wrap(err, "Can not save highest sources, targets and proposals for account")
		}
	}
	return nil
}

// SaveHighestData save the highest source, target and proposal for account
func SaveHighestData(acc core.ValidatorAccount, store *inmemory.InMemStore, accountFlags *CreateAccountFlagValues, index int) error {
	highestIndex := index
	if !accountFlags.accumulate && len(accountFlags.privateKeys) <= 1 {
		highestIndex = 0
	}

	// add minimal attestation protection data
	minimalAtt := &phase0.AttestationData{
		Source: &phase0.Checkpoint{Epoch: phase0.Epoch(accountFlags.highestSources[highestIndex])},
		Target: &phase0.Checkpoint{Epoch: phase0.Epoch(accountFlags.highestTargets[highestIndex])},
	}
	if err := store.SaveHighestAttestation(acc.ValidatorPublicKey(), minimalAtt); err != nil {
		return errors.Wrap(err, "failed to save highest attestation")
	}

	// add minimal proposal protection data
	if err := store.SaveHighestProposal(acc.ValidatorPublicKey(), phase0.Slot(accountFlags.highestProposals[highestIndex])); err != nil {
		return errors.Wrap(err, "failed to save highest proposal")
	}
	return nil
}

// ValidateHighestValues Performs basic validation for account highest attestation/proposal values
func ValidateHighestValues(accountFlagValues CreateAccountFlagValues) error {
	if len(accountFlagValues.privateKeys) > 0 {
		errorExplain := "length for seedless accounts need to be equal to private keys count"
		privateKeysCount := len(accountFlagValues.privateKeys)

		if len(accountFlagValues.highestSources) != privateKeysCount {
			return errors.Errorf("highest sources %v", errorExplain)
		}
		if len(accountFlagValues.highestTargets) != privateKeysCount {
			return errors.Errorf("highest targets %v", errorExplain)
		}
		if len(accountFlagValues.highestProposals) != privateKeysCount {
			return errors.Errorf("highest proposals %v", errorExplain)
		}
	} else if accountFlagValues.accumulate {
		if len(accountFlagValues.highestSources) != (accountFlagValues.index + 1) {
			return errors.Errorf("highest sources length when the accumulate flag is true need to be equal to index")
		}
		if len(accountFlagValues.highestTargets) != (accountFlagValues.index + 1) {
			return errors.Errorf("highest targets length when the accumulate flag is true need to be index")
		}
		if len(accountFlagValues.highestProposals) != (accountFlagValues.index + 1) {
			return errors.Errorf("highest proposals length when the accumulate flag is true need to be index")
		}
	} else {
		if len(accountFlagValues.highestSources) != 1 {
			return errors.Errorf("highest sources length when the accumulate flag is false need to be 1")
		}
		if len(accountFlagValues.highestTargets) != 1 {
			return errors.Errorf("highest targets length when the accumulate flag is false need to be 1")
		}
		if len(accountFlagValues.highestProposals) != 1 {
			return errors.Errorf("highest proposals length when the accumulate flag is false need to be 1")
		}
	}
	return nil
}
