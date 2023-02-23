package handler

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
)

// CredentialsFlagValues keeps all collected values for seed
type CredentialsFlagValues struct {
	index      int
	seedBytes  []byte
	accumulate bool
	validators []*core.ValidatorInfo
	network    core.Network
}

// Credentials creates a new wallet account(s) and prints the storage.
func (h *Account) Credentials(cmd *cobra.Command, args []string) error {
	err := core.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	credentialsFlags, err := CollectCredentialsFlags(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to collect credentials flags")
	}

	// Initialize store
	store := inmemory.NewInMemStore(credentialsFlags.network)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	// Create new key vault
	_, err = eth2keymanager.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	signedBLSToExecutionChanges := make([]*capella.SignedBLSToExecutionChange, 0)

	for i := 0; i <= credentialsFlags.index; i++ {
		var index int
		if credentialsFlags.accumulate {
			index = i
		} else {
			index = credentialsFlags.index
		}

		signedBLSToExecutionChange, err := wallet.CreateSignedBLSToExecutionChange(credentialsFlags.validators[i], credentialsFlags.seedBytes, &index)
		if err != nil {
			return errors.Wrap(err, "failed to build BLS to execution change")
		}
		signedBLSToExecutionChanges = append(signedBLSToExecutionChanges, signedBLSToExecutionChange)

		if !credentialsFlags.accumulate {
			break
		}
	}

	err = h.printer.JSON(signedBLSToExecutionChanges)
	if err != nil {
		return errors.Wrap(err, "failed to print signedBLSToExecutionChanges JSON")
	}

	return nil
}

// CollectCredentialsFlags returns collected flags for seed
func CollectCredentialsFlags(cmd *cobra.Command) (*CredentialsFlagValues, error) {
	credentialsFlagValues := CredentialsFlagValues{}

	// Get seed flag value.
	seedFlagValue, err := flag.GetSeedFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	// Get seed bytes
	seedBytes, err := hex.DecodeString(seedFlagValue)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode seed")
	}
	credentialsFlagValues.seedBytes = seedBytes

	// Get accumulate flag value.
	accumulateFlagValue, err := rootcmd.GetAccumulateFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the accumulate flag value")
	}
	credentialsFlagValues.accumulate = accumulateFlagValue

	// Get index flag value.
	indexFlagValue, err := flag.GetIndexFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the index flag value")
	}
	credentialsFlagValues.index = indexFlagValue

	// Get network flag value.
	network, err := rootcmd.GetNetworkFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the network flag value")
	}
	credentialsFlagValues.network = network

	// Get validators info flag value.
	validators, err := flag.GetValidatorInfoFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the validators info flag value")
	}
	credentialsFlagValues.validators = validators

	return &credentialsFlagValues, nil
}
