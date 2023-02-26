package handler

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/signer"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
)

// VoluntaryExitFlagValues keeps all collected values for seed
type VoluntaryExitFlagValues struct {
	index              int
	seedBytes          []byte
	accumulate         bool
	currentForkVersion phase0.Version
	epoch              int
	validators         []*core.ValidatorInfo
	network            core.Network
}

// VoluntaryExit creates a new wallet account(s) and prints the storage.
func (h *Account) VoluntaryExit(cmd *cobra.Command, args []string) error {
	err := core.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	voluntaryExitFlags, err := CollectVoluntaryExitFlags(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to collect voluntary exit flags")
	}

	// Initialize store
	store := inmemory.NewInMemStore(voluntaryExitFlags.network)
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

	// Compute domain
	genesisValidatorsRoot := store.Network().GenesisValidatorsRoot()
	domainBytes, err := types.ComputeDomain(types.DomainVoluntaryExit, voluntaryExitFlags.currentForkVersion[:], genesisValidatorsRoot[:])
	if err != nil {
		return errors.Wrap(err, "failed to calculate domain")
	}
	var domain phase0.Domain
	copy(domain[:], domainBytes)

	simpleSigner := signer.NewSimpleSigner(wallet, nil, store.Network())
	signedVoluntaryExits := make([]*phase0.SignedVoluntaryExit, 0)

	for i := 0; i <= voluntaryExitFlags.index; i++ {
		var index int
		if voluntaryExitFlags.accumulate {
			index = i
		} else {
			index = voluntaryExitFlags.index
		}

		acc, err := wallet.CreateValidatorAccount(voluntaryExitFlags.seedBytes, &index)
		if err != nil {
			return errors.Wrap(err, "failed to create validator account")
		}
		voluntaryExit := &phase0.VoluntaryExit{
			Epoch:          phase0.Epoch(voluntaryExitFlags.epoch),
			ValidatorIndex: voluntaryExitFlags.validators[i].Index,
		}
		signature, _, err := simpleSigner.SignVoluntaryExit(voluntaryExit, domain, acc.ValidatorPublicKey())
		if err != nil {
			return errors.Wrap(err, "failed to sign voluntary exit")
		}

		signedVoluntaryExit := &phase0.SignedVoluntaryExit{
			Message: voluntaryExit,
		}
		copy(signedVoluntaryExit.Signature[:], signature)
		signedVoluntaryExits = append(signedVoluntaryExits, signedVoluntaryExit)

		if !voluntaryExitFlags.accumulate {
			break
		}
	}

	err = h.printer.JSON(signedVoluntaryExits)
	if err != nil {
		return errors.Wrap(err, "failed to print signed voluntary exit JSON")
	}

	return nil
}

// CollectVoluntaryExitFlags returns collected flags for seed
func CollectVoluntaryExitFlags(cmd *cobra.Command) (*VoluntaryExitFlagValues, error) {
	voluntaryExitFlagValues := VoluntaryExitFlagValues{}

	// Get network flag value.
	network, err := rootcmd.GetNetworkFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the network flag value")
	}
	voluntaryExitFlagValues.network = network

	// Get seed flag value.
	seedFlagValue, err := rootcmd.GetSeedFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	// Get seed bytes
	seedBytes, err := hex.DecodeString(seedFlagValue)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode seed")
	}
	voluntaryExitFlagValues.seedBytes = seedBytes

	// Get accumulate flag value.
	accumulateFlagValue, err := rootcmd.GetAccumulateFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the accumulate flag value")
	}
	voluntaryExitFlagValues.accumulate = accumulateFlagValue

	// Get index flag value.
	indexFlagValue, err := rootcmd.GetIndexFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the index flag value")
	}
	voluntaryExitFlagValues.index = indexFlagValue

	// Get current fork version flag value.
	currentForkVersionFlagValue, err := flag.GetCurrentForkVersionFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the current fork version flag value")
	}
	voluntaryExitFlagValues.currentForkVersion = currentForkVersionFlagValue

	// Get epoch flag value.
	epochFlagValue, err := flag.GetEpochFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the index flag value")
	}
	voluntaryExitFlagValues.epoch = epochFlagValue

	// Get validators info flag value.
	validators, err := flag.GetVoluntaryExitInfoFlagValue(cmd)
	if err != nil {
		return nil, err
	}
	voluntaryExitFlagValues.validators = validators

	return &voluntaryExitFlagValues, nil
}
