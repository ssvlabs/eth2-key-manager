package flag

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ssvlabs/eth2-key-manager/cli/util/cliflag"
)

// Flag names.
const (
	privateKeyFlag       = "private-key"
	highestKnownSource   = "highest-source"
	highestKnownTarget   = "highest-target"
	highestKnownProposal = "highest-proposal"
)

// AddPrivateKeyFlag adds private key to the command
func AddPrivateKeyFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, privateKeyFlag, "", "private key", false)
}

// GetPrivateKeyFlagValue returns the value of private key
func GetPrivateKeyFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(privateKeyFlag)
}

// GetPrivateKeyFlagName gets privateKeyFlag name
func GetPrivateKeyFlagName() string {
	return privateKeyFlag
}

// AddHighestSourceFlag adds the highest source flag to the command
func AddHighestSourceFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, highestKnownSource, "", "Array of highest known sources for an array of validators", true)
}

// GetHighestSourceFlagValue gets the highest source flag from the command
func GetHighestSourceFlagValue(c *cobra.Command) ([]uint64, error) {
	str, err := c.Flags().GetString(highestKnownSource)
	if err != nil {
		return nil, err
	}
	return stringSliceToUint64Slice(str)
}

// AddHighestTargetFlag adds the highest target flag to the command
func AddHighestTargetFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, highestKnownTarget, "", "Array of highest known targets for an array of validators", true)
}

// GetHighestTargetFlagValue gets the highest target flag from the command
func GetHighestTargetFlagValue(c *cobra.Command) ([]uint64, error) {
	str, err := c.Flags().GetString(highestKnownTarget)
	if err != nil {
		return nil, err
	}
	return stringSliceToUint64Slice(str)
}

// AddHighestProposalFlag adds the highest proposal flag to the command
func AddHighestProposalFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, highestKnownProposal, "", "Array of highest known proposed blocks (slot) for an array of validators", true)
}

// GetHighestProposalFlagValue gets the highest proposal flag from the command
func GetHighestProposalFlagValue(c *cobra.Command) ([]uint64, error) {
	str, err := c.Flags().GetString(highestKnownProposal)
	if err != nil {
		return nil, err
	}
	return stringSliceToUint64Slice(str)
}

func stringSliceToUint64Slice(str string) ([]uint64, error) {
	strs := strings.Split(str, ",")
	ret := make([]uint64, len(strs))
	for i, s := range strs {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ret[i] = n
	}
	return ret, nil
}
