package flag

import (
	"encoding/hex"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Flag names.
const (
	validatorPublicKey     = "validator-public-key"
	validatorIndex         = "validator-index"
	epochFlag              = "epoch"
	currentForkVersionFlag = "current-fork-version"
)

// AddEpochFlag adds the epoch flag to the command
func AddEpochFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, epochFlag, 0, "epoch", true)
}

// GetEpochFlagValue gets the epoch flag from the command
func GetEpochFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(epochFlag)
}

// AddCurrentForkVersionFlag adds the current fork version flag to the command
func AddCurrentForkVersionFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, currentForkVersionFlag, "", "current fork version (ForkVersionLength: 4 bytes)", true)
}

// GetCurrentForkVersionFlagValue gets the current fork version flag from the command
func GetCurrentForkVersionFlagValue(c *cobra.Command) (phase0.Version, error) {
	var currentForkVersion phase0.Version
	currentForkVersionFlagValue, err := c.Flags().GetString(currentForkVersionFlag)
	if err != nil {
		return currentForkVersion, err
	}
	version, err := hex.DecodeString(strings.TrimPrefix(currentForkVersionFlagValue, "0x"))
	if err != nil {
		return currentForkVersion, errors.Wrap(err, "invalid current fork version supplied")
	}
	if len(version) != phase0.ForkVersionLength {
		return currentForkVersion, errors.New("invalid length for current fork version")
	}

	copy(currentForkVersion[:], version)

	return currentForkVersion, nil
}

// AddValidatorPublicKeyFlag adds the validator public key flag to the command
func AddValidatorPublicKeyFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, validatorPublicKey, "", "validator public key", true)
}

// GetValidatorPublicKeyFlagValue gets the validator public key flag from the command
func GetValidatorPublicKeyFlagValue(c *cobra.Command) (phase0.BLSPubKey, error) {
	var validatorBLSPubKey phase0.BLSPubKey
	validatorPublicKeyValue, err := c.Flags().GetString(validatorPublicKey)
	if err != nil {
		return validatorBLSPubKey, errors.Wrap(err, "failed to retrieve the validator public key flag value")
	}

	validatorPubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(validatorPublicKeyValue, "0x"))
	if err != nil {
		return validatorBLSPubKey, errors.Wrap(err, "invalid validator public key supplied")
	}
	if len(validatorPubKeyBytes) != phase0.PublicKeyLength {
		return validatorBLSPubKey, errors.New("invalid length for validator public key")
	}
	copy(validatorBLSPubKey[:], validatorPubKeyBytes)

	return validatorBLSPubKey, nil
}

// AddValidatorIndexFlag adds the validator index flag to the command
func AddValidatorIndexFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, validatorIndex, 0, "validator index", true)
}

// GetValidatorIndexFlagValue gets the validator index flag from the command
func GetValidatorIndexFlagValue(c *cobra.Command) (phase0.ValidatorIndex, error) {
	str, err := c.Flags().GetInt(validatorIndex)
	if err != nil {
		return 0, err
	}

	return phase0.ValidatorIndex(str), nil
}

// GetVoluntaryExitInfoFlagValue gets the voluntary exit info flag from the command
func GetVoluntaryExitInfoFlagValue(c *cobra.Command) (*core.ValidatorInfo, error) {
	validatorIndex, err := GetValidatorIndexFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse validator index")
	}

	validatorPubKey, err := GetValidatorPublicKeyFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse validator public key")
	}

	validatorInfo := &core.ValidatorInfo{
		Index:  validatorIndex,
		Pubkey: validatorPubKey,
	}

	return validatorInfo, nil
}
