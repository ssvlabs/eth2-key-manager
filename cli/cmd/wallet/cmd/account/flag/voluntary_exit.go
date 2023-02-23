package flag

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/publickey/flag"
	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Flag names.
const (
	epochFlag = "epoch"
)

// AddEpochFlag adds the epoch flag to the command
func AddEpochFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, epochFlag, 0, "epoch", true)
}

// GetEpochFlagValue gets the epoch flag from the command
func GetEpochFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(epochFlag)
}

// GetVoluntaryExitInfoFlagValue gets the voluntary exit info flag from the command
func GetVoluntaryExitInfoFlagValue(c *cobra.Command) ([]*core.ValidatorInfo, error) {
	validatorIndices, err := GetValidatorIndexFlagValue(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse validator indices")
	}

	indexFlagValue, err := flag.GetIndexFlagValue(c)
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
			Index: phase0.ValidatorIndex(validatorIndices[i]),
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
