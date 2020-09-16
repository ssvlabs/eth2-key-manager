package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Flag names.
const (
	networkFlag = "network"
)

// AddNetworkFlag adds the network flag to the command
func AddNetworkFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, networkFlag, string(core.TestNetwork), "Ethereum network", false)
}

// GetNetworkFlagValue gets the network flag from the command
func GetNetworkFlagValue(c *cobra.Command) (core.Network, error) {
	networkValue, err := c.Flags().GetString(networkFlag)
	if err != nil {
		return "", err
	}

	switch networkValue {
	case string(core.MainNetwork):
		return core.MainNetwork, nil
	case string(core.LaunchTestNetwork):
		return core.LaunchTestNetwork, nil
	case string(core.TestNetwork):
		return core.TestNetwork, nil
	default:
		return "", errors.New("undefined network")
	}
}
