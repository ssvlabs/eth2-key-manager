package handler

import (
	"bytes"
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/signer"
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

var domainBlsToExecutionChange = types.DomainType{0x0a, 0x00, 0x00, 0x00}

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

	// set the context for wallet to use withdrawal key as primary key
	wallet.SetContext(&core.WalletContext{
		Storage:        store,
		WithdrawalMode: true,
	})

	// Compute domain
	genesisValidatorsRoot, err := store.Network().GenesisValidatorsRoot()
	if err != nil {
		return errors.Wrap(err, "failed to get genesis validators root")
	}
	genesisForkVersion := store.Network().GenesisForkVersion()
	domainBytes, err := types.ComputeDomain(domainBlsToExecutionChange, genesisForkVersion[:], genesisValidatorsRoot[:])
	if err != nil {
		return errors.Wrap(err, "failed to calculate domain")
	}
	var domain phase0.Domain
	copy(domain[:], domainBytes)

	simpleSigner := signer.NewSimpleSigner(wallet, nil, store.Network())
	signedBLSToExecutionChanges := make([]*capella.SignedBLSToExecutionChange, 0)

	for i := 0; i <= credentialsFlags.index; i++ {
		var index int
		if credentialsFlags.accumulate {
			index = i
		} else {
			index = credentialsFlags.index
		}
		validator := credentialsFlags.validators[i]

		// Its actually withdrawal account, since the wallet is in withdrawal mode
		acc, err := wallet.CreateValidatorAccount(credentialsFlags.seedBytes, &index)
		if err != nil {
			return errors.Wrap(err, "failed to create withdrawal account")
		}

		// Since the wallet is in withdrawal mode the validator account is the withdrawal account
		derivedWithdrawalPubKey := acc.ValidatorPublicKey()
		derivedValidatorPubKey := acc.WithdrawalPublicKey()

		// validation that the derived validation public key is the same as the one in the validator info
		if !bytes.Equal(derivedValidatorPubKey, validator.Pubkey[:]) {
			derivedPubKey := "0x" + hex.EncodeToString(derivedValidatorPubKey)
			providedPubKey := validator.Pubkey.String()
			return errors.Errorf("derived validator public key: %s, does not match with the provided one: %s", derivedPubKey, providedPubKey)
		}

		// validation that the derived withdrawal credentials are the same as the one in the validator info
		withdrawalCredentials := util.SHA256(derivedWithdrawalPubKey)
		withdrawalCredentials[0] = byte(0) // BLS_WITHDRAWAL_PREFIX
		if !bytes.Equal(withdrawalCredentials, validator.WithdrawalCredentials) {
			derivedCreds := "0x" + hex.EncodeToString(withdrawalCredentials)
			providedCreds := "0x" + hex.EncodeToString(validator.WithdrawalCredentials)
			return errors.Errorf("derived withdrawal credentials: %s, does not match with the provided one: %s", derivedCreds, providedCreds)
		}

		blsToExecutionChange := &capella.BLSToExecutionChange{
			ValidatorIndex:     validator.Index,
			ToExecutionAddress: validator.ToExecutionAddress,
		}
		copy(blsToExecutionChange.FromBLSPubkey[:], derivedWithdrawalPubKey)

		signature, _, err := simpleSigner.SignBLSToExecutionChange(blsToExecutionChange, domain, derivedWithdrawalPubKey)
		if err != nil {
			return errors.Wrap(err, "failed to sign voluntary exit")
		}

		signedBLSToExecutionChange := &capella.SignedBLSToExecutionChange{
			Message: blsToExecutionChange,
		}
		copy(signedBLSToExecutionChange.Signature[:], signature)
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
	seedFlagValue, err := rootcmd.GetSeedFlagValue(cmd)
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
	indexFlagValue, err := rootcmd.GetIndexFlagValue(cmd)
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
		return nil, err
	}
	credentialsFlagValues.validators = validators

	return &credentialsFlagValues, nil
}
