package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Flag names.
const (
	networkFlag    = "network"
	accumulateFlag = "accumulate"
)

// AddNetworkFlag adds the network flag to the command
func AddNetworkFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, networkFlag, "", "Ethereum network", false)
}

// GetNetworkFlagValue gets the network flag from the command
func GetNetworkFlagValue(c *cobra.Command) (core.Network, error) {
	networkValue, err := c.Flags().GetString(networkFlag)
	if err != nil {
		return "", err
	}

	ret := core.NetworkFromString(networkValue)
	if len(ret) == 0 {
		return "", errors.New("unknown network")
	}

	return ret, nil
}

// AddAccumulateFlag adds the accumulate flag to the command
func AddAccumulateFlag(c *cobra.Command) {
	cliflag.AddPersistentBoolFlag(c, accumulateFlag, false, "accumulate accounts", false)
}

// GetAccumulateFlagValue gets the accumulate flag from the command
func GetAccumulateFlagValue(c *cobra.Command) (bool, error) {
	return c.Flags().GetBool(accumulateFlag)
}
