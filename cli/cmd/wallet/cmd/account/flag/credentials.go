package flag

import (
	"encoding/hex"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Flag names.
const (
	validatorIndex        = "validator-index"
	validatorPublicKey    = "validator-public-key"
	withdrawalCredentials = "withdrawal-credentials"
	toExecutionAddress    = "to-execution-address"
)

// AddValidatorIndexFlag adds the validator index flag to the command
func AddValidatorIndexFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, validatorIndex, "", "comma separate string of validator indices", true)
}

// GetValidatorIndexFlagValue gets the validator index flag from the command
func GetValidatorIndexFlagValue(c *cobra.Command) ([]uint64, error) {
	str, err := c.Flags().GetString(validatorIndex)
	if err != nil {
		return nil, err
	}
	return stringSliceToUint64Slice(str)
}

// AddValidatorPublicKeyFlag adds the validator public key flag to the command
func AddValidatorPublicKeyFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, validatorPublicKey, "", "comma separate string of validator public keys", true)
}

// GetValidatorPublicKeyFlagValue gets the validator index flag from the command
func GetValidatorPublicKeyFlagValue(c *cobra.Command) ([]phase0.BLSPubKey, error) {
	validatorPublicKeyValues, err := c.Flags().GetString(validatorPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the validator public key flag value")
	}

	validatorPubKeys := strings.Split(validatorPublicKeyValues, ",")
	validatorBLSPubKeys := make([]phase0.BLSPubKey, 0)
	for _, pk := range validatorPubKeys {
		validatorPubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(pk, "0x"))
		if err != nil {
			return nil, errors.Wrap(err, "invalid validator public key supplied")
		}
		if len(validatorPubKeyBytes) != phase0.PublicKeyLength {
			return nil, errors.New("invalid length for validator public key")
		}

		var validatorBLSPubKey phase0.BLSPubKey
		copy(validatorBLSPubKey[:], validatorPubKeyBytes)
		validatorBLSPubKeys = append(validatorBLSPubKeys, validatorBLSPubKey)
	}
	return validatorBLSPubKeys, nil
}

// AddWithdrawalCredentialsFlag adds withdrawal credentials to the command
func AddWithdrawalCredentialsFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, withdrawalCredentials, "", "comma separate string of withdrawal credentials", true)
}

// GetWithdrawalCredentialsFlagValue returns the value of withdrawal credentials
func GetWithdrawalCredentialsFlagValue(c *cobra.Command) ([][]byte, error) {
	withdrawalCredentialsValues, err := c.Flags().GetString(withdrawalCredentials)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the withdrawal credentials flag value")
	}

	withdrawalCreds := strings.Split(withdrawalCredentialsValues, ",")
	withdrawalCredentialsList := make([][]byte, 0)
	for _, cred := range withdrawalCreds {
		withdrawalCredsBytes, err := hex.DecodeString(strings.TrimPrefix(cred, "0x"))
		if err != nil {
			return nil, errors.Wrap(err, "invalid withdrawal credentials supplied")
		}
		if len(withdrawalCredsBytes) != 32 {
			return nil, errors.New("invalid length for withdrawal credentials")
		}

		if withdrawalCredsBytes[0] != byte(0) {
			return nil, errors.New("non-BLS withdrawal credentials supplied")
		}
		withdrawalCredentialsList = append(withdrawalCredentialsList, withdrawalCredsBytes)
	}
	return withdrawalCredentialsList, nil
}

// AddToExecutionAddressFlag adds the validator execution withdrawal address flag to the command
func AddToExecutionAddressFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, toExecutionAddress, "", "comma separate string of to execution addresses", true)
}

// GetToExecutionAddressFlagValue gets the validator withdrawal address flag from the command
func GetToExecutionAddressFlagValue(c *cobra.Command) ([]bellatrix.ExecutionAddress, error) {
	toExecutionAddressValues, err := c.Flags().GetString(toExecutionAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the withdrawal address flag value")
	}

	toExecutionAddresses := strings.Split(toExecutionAddressValues, ",")
	toExecutionAddressList := make([]bellatrix.ExecutionAddress, 0)
	for _, wa := range toExecutionAddresses {
		toExecutionAddressBytes, err := hex.DecodeString(strings.TrimPrefix(wa, "0x"))
		if err != nil {
			return nil, errors.Wrap(err, "invalid to execution address supplied")
		}
		if len(toExecutionAddressBytes) != bellatrix.ExecutionAddressLength {
			return nil, errors.New("invalid length for to execution address")
		}

		var toExecutionAdd bellatrix.ExecutionAddress
		copy(toExecutionAdd[:], toExecutionAddressBytes)
		toExecutionAddressList = append(toExecutionAddressList, toExecutionAdd)
	}
	return toExecutionAddressList, nil
}

// GetValidatorInfoFlagValue gets the validator info flag from the command
func GetValidatorInfoFlagValue(c *cobra.Command) ([]*core.ValidatorInfo, error) {
	validatorIndices, err := GetValidatorIndexFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse validator indices")
	}

	validatorPubKeys, err := GetValidatorPublicKeyFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse validator public keys")
	}

	validatorWithdrawalCredentials, err := GetWithdrawalCredentialsFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse withdrawal credentials")
	}

	toExecutionAddresses, err := GetToExecutionAddressFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse to execution addresses")
	}

	if len(validatorIndices) != len(validatorPubKeys) || len(validatorPubKeys) != len(validatorWithdrawalCredentials) || len(validatorWithdrawalCredentials) != len(toExecutionAddresses) {
		return nil, errors.New("validator indices, public keys, withdrawal credentials and to execution addresses must be of equal length")
	}

	indexFlagValue, err := rootcmd.GetIndexFlagValue(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the index flag value")
	}

	accumulate, err := rootcmd.GetAccumulateFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse accumulate flag")
	}

	if accumulate && indexFlagValue+1 != len(validatorIndices) {
		return nil, errors.New("index flag value must be one less than the number of validator indices")
	}

	if !accumulate && len(validatorIndices) > 1 {
		return nil, errors.New("only one validator can be specified if accumulate is false")
	}

	validatorInfoList := make([]*core.ValidatorInfo, len(validatorIndices))
	for i := 0; i < len(validatorIndices); i++ {
		validatorInfoList[i] = &core.ValidatorInfo{
			Index:                 phase0.ValidatorIndex(validatorIndices[i]),
			Pubkey:                validatorPubKeys[i],
			WithdrawalCredentials: validatorWithdrawalCredentials[i],
			ToExecutionAddress:    toExecutionAddresses[i],
		}
	}

	if accumulate && indexFlagValue+1 != len(validatorInfoList) {
		return nil, errors.New("index flag value must be one less than the number of validator info")
	}

	if !accumulate && len(validatorInfoList) > 1 {
		return nil, errors.New("only one validator info can be structured if accumulate is false")
	}

	return validatorInfoList, nil
}
