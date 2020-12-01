package handler

import (
	"encoding/hex"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Account creates a new wallet account and prints the storage.
func (h *Account) Create(cmd *cobra.Command, args []string) error {
	err := core.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	// Get index flag.
	indexFlagValue, err := flag.GetIndexFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the index flag value")
	}

	// Get seed flag.
	seedFlagValue, err := flag.GetSeedFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	seedBytes, err := hex.DecodeString(seedFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode seed")
	}

	// Get accumulate flag.
	accumulateFlagValue, err := flag.GetAccumulateFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the accumulate flag value")
	}

	// Get response-type flag.
	responseType, err := flag.GetResponseTypeFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the response type value")
	}

	// Get minimals slashing data flag
	highestSources, err := flag.GetHighestSourceFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}
	highestTargets, err := flag.GetHighestTargetFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the minimal slashing data value")
	}

	if accumulateFlagValue {
		if len(highestSources) != (indexFlagValue + 1) {
			return errors.Errorf("highest sources length when the accumulate flag is true need to be equal to indexFlagValue")
		}
		if len(highestTargets) != (indexFlagValue + 1) {
			return errors.Errorf("highest targets length when the accumulate flag is true need to be indexFlagValue")
		}
	} else {
		if len(highestSources) != 1 {
			return errors.Errorf("highest sources length when the accumulate flag is false need to be 1")
		}
		if len(highestTargets) != 1 {
			return errors.Errorf("highest targets length when the accumulate flag is false need to be 1")
		}
	}

	// TODO get rid of network
	store := in_memory.NewInMemStore(core.PyrmontNetwork)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	_, err = eth2keymanager.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	if accumulateFlagValue {
		for i := 0; i <= indexFlagValue; i++ {
			acc, err := wallet.CreateValidatorAccount(seedBytes, &i)
			if err != nil {
				return errors.Wrap(err, "failed to create validator account")
			}

			minimalAtt := &eth.AttestationData{
				Source: &eth.Checkpoint{Epoch: uint64(highestSources[i])},
				Target: &eth.Checkpoint{Epoch: uint64(highestTargets[i])},
			}

			if err := store.SaveHighestAttestation(acc.ValidatorPublicKey(), minimalAtt); err != nil {
				return errors.Wrap(err, "failed to set validator minimal slashing protection")
			}
		}
	} else {
		acc, err := wallet.CreateValidatorAccount(seedBytes, &indexFlagValue)
		if err != nil {
			return errors.Wrap(err, "failed to create validator account")
		}

		minimalAtt := &eth.AttestationData{
			Source: &eth.Checkpoint{Epoch: uint64(highestSources[0])},
			Target: &eth.Checkpoint{Epoch: uint64(highestTargets[0])},
		}
		if err := store.SaveHighestAttestation(acc.ValidatorPublicKey(), minimalAtt); err != nil {
			return errors.Wrap(err, "failed to set validator minimal slashing protection")
		}
	}

	if responseType == flag.StorageResponseType {
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
		accObj := map[string]string{
			"id":               a.ID().String(),
			"name":             a.Name(),
			"validationPubKey": hex.EncodeToString(a.ValidatorPublicKey()),
			"withdrawalPubKey": hex.EncodeToString(a.WithdrawalPublicKey()),
		}
		accounts = append(accounts, accObj)
	}

	if accumulateFlagValue {
		err = h.printer.JSON(accounts)
	} else if len(accounts) > 0 {
		err = h.printer.JSON(accounts[0])
	}
	if err != nil {
		return errors.Wrap(err, "failed to print accounts JSON")
	}
	return nil
}
