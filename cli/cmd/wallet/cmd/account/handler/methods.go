package handler

import (
	"encoding/hex"
	types "github.com/prysmaticlabs/eth2-types"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cmd2 "github.com/bloxapp/eth2-key-manager/cli/cmd"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// CreateAccountFlagValues keeps all collected values for seed and seedless modes
type CreateAccountFlagValues struct {
	index            int
	indexFrom        int
	seed             string
	seedBytes        []byte
	privateKey       [][]byte
	accumulate       bool
	responseType     flag.ResponseType
	highestSources   []uint64
	highestTargets   []uint64
	highestProposals []uint64
	network          core.Network
}

// CheckHighestValues Performs basic checks for account flags
func CheckHighestValues(accountFlagValues CreateAccountFlagValues) error {
	if len(accountFlagValues.privateKey) > 0 {
		errorExplain := "length for seedless accounts need to be equal to private keys count"
		privateKeysCount := len(accountFlagValues.privateKey)

		if len(accountFlagValues.highestSources) != privateKeysCount {
			return errors.Errorf("highest sources " + errorExplain)
		}
		if len(accountFlagValues.highestTargets) != privateKeysCount {
			return errors.Errorf("highest targets " + errorExplain)
		}
		if len(accountFlagValues.highestProposals) != privateKeysCount {
			return errors.Errorf("highest proposals " + errorExplain)
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

// CollectAccountFlags returns collected flags for seed and seedless modes
func CollectAccountFlags(cmd *cobra.Command, seedless bool) (*CreateAccountFlagValues, error) {
	accountFlagValues := CreateAccountFlagValues{}

	// Seedless mode
	if seedless == true {
		// Private key
		indexFromValue, err := flag.GetIndexFromFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the index-from flag value")
		}
		accountFlagValues.indexFrom = indexFromValue

		// Private key
		privateKeyValues, err := flag.GetPrivateKeyFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the private key flag value")
		}

		privateKeys := strings.Split(privateKeyValues, ",")
		for i := 0; i < len(privateKeys); i++ {
			privateKeyBytes, err := hex.DecodeString(privateKeys[i])
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert private key string to bytes")
			}
			accountFlagValues.privateKey = append(accountFlagValues.privateKey, privateKeyBytes)
		}
		// Seed mode
	} else {
		// Get index flag.
		indexFlagValue, err := flag.GetIndexFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the index flag value")
		}
		accountFlagValues.index = indexFlagValue

		// Seed flag
		seedFlagValue, err := flag.GetSeedFlagValue(cmd)
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

		// Get accumulate flag.
		accumulateFlagValue, err := flag.GetAccumulateFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the accumulate flag value")
		}
		accountFlagValues.accumulate = accumulateFlagValue
	}

	// Get response-type flag.
	responseType, err := flag.GetResponseTypeFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the response type value")
	}
	accountFlagValues.responseType = responseType

	// Get minimals slashing data flag
	highestSources, err := flag.GetHighestSourceFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}
	accountFlagValues.highestSources = highestSources

	highestTargets, err := flag.GetHighestTargetFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}
	accountFlagValues.highestTargets = highestTargets

	highestProposals, err := flag.GetHighestProposalFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}
	accountFlagValues.highestProposals = highestProposals

	// Check highest values
	highestValuesError := CheckHighestValues(accountFlagValues)
	if highestValuesError != nil {
		return nil, highestValuesError
	}

	// Get network value
	network, err := cmd2.GetNetworkFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to network")
	}
	accountFlagValues.network = network

	return &accountFlagValues, nil
}

// GenerateAccounts generates account by index using provided account flags
func GenerateAccounts(wallet core.Wallet, store *inmemory.InMemStore, index int, accountFlags *CreateAccountFlagValues) error {
	var acc core.ValidatorAccount
	var err error

	if len(accountFlags.privateKey) > 0 {
		for i := accountFlags.indexFrom; i < accountFlags.indexFrom+len(accountFlags.privateKey); i++ {
			acc, err = wallet.CreateValidatorAccountFromPrivateKey(accountFlags.privateKey[i-accountFlags.indexFrom], &i)
			if err != nil {
				return errors.Wrap(err, "failed to create validator account")
			}

			err = SaveHighestData(acc, store, accountFlags, i-accountFlags.indexFrom)
			if err != nil {
				return errors.Wrap(err, "Can not save highest sources, targets and proposals for account")
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

// SaveHighestData save the highest sources, targets and proposals for account
func SaveHighestData(acc core.ValidatorAccount, store *inmemory.InMemStore, accountFlags *CreateAccountFlagValues, index int) error {
	highestIndex := index
	if accountFlags.accumulate != true && len(accountFlags.privateKey) <= 1 {
		highestIndex = 0
	}

	// add minimal attestation protection data
	minimalAtt := &eth.AttestationData{
		Source: &eth.Checkpoint{Epoch: types.Epoch(accountFlags.highestSources[highestIndex])},
		Target: &eth.Checkpoint{Epoch: types.Epoch(accountFlags.highestTargets[highestIndex])},
	}
	if err := store.SaveHighestAttestation(acc.ValidatorPublicKey(), minimalAtt); err != nil {
		return errors.Wrap(err, "failed to set validator minimal slashing protection")
	}

	// add minimal proposal protection data
	minimalProposal := &eth.BeaconBlock{
		Slot: types.Slot(accountFlags.highestProposals[highestIndex]),
	}
	if err := store.SaveHighestProposal(acc.ValidatorPublicKey(), minimalProposal); err != nil {
		return errors.Wrap(err, "failed to set validator minimal slashing protection")
	}
	return nil
}

// BuildAndPrintAccounts builds all accounts or one account depending on seedless flag
func (h *Account) BuildAndPrintAccounts(accountFlags *CreateAccountFlagValues) error {
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

	if accountFlags.responseType == flag.StorageResponseType {
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
		if len(accountFlags.privateKey) > 0 {
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

	if accountFlags.accumulate || len(accountFlags.privateKey) > 1 {
		err = h.printer.JSON(accounts)
	} else if len(accounts) > 0 {
		err = h.printer.JSON(accounts[0])
	}
	if err != nil {
		return errors.Wrap(err, "failed to print accounts JSON")
	}
	return nil
}
