package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
	"github.com/bloxapp/eth2-key-manager/core"
)

// ResponseType represents the network.
type ResponseType string

// Available response types.
const (
	// StorageResponseType represents the storage response type.
	StorageResponseType ResponseType = "storage"

	// ObjectResponseType represents the storage response type.
	ObjectResponseType ResponseType = "object"
)

// ResponseTypeFromString returns response type from the given string value
func ResponseTypeFromString(n string) ResponseType {
	switch n {
	case string(StorageResponseType):
		return StorageResponseType
	case string(ObjectResponseType):
		return ObjectResponseType
	default:
		panic(fmt.Sprintf("undefined response type %s", n))
	}
}

// Flag names.
const (
	networkFlag      = "network"
	accumulateFlag   = "accumulate"
	seedFlag         = "seed"
	indexFlag        = "index"
	responseTypeFlag = "response-type"
)

// AddNetworkFlag adds the network flag to the command
func AddNetworkFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, networkFlag, "", "Ethereum network", true)
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

// AddSeedFlag adds the seed flag to the command
func AddSeedFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, seedFlag, "", "seed", false)
}

// GetSeedFlagValue gets the seed flag from the command
func GetSeedFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(seedFlag)
}

// AddIndexFlag adds the index flag to the command
func AddIndexFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, indexFlag, 0, "public key index", true)
}

// GetIndexFlagValue gets the index flag from the command
func GetIndexFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(indexFlag)
}

// AddResponseTypeFlag adds the response-type flag to the command
func AddResponseTypeFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, responseTypeFlag, string(StorageResponseType), "response type", false)
}

// GetResponseTypeFlagValue gets the response-type flag from the command
func GetResponseTypeFlagValue(c *cobra.Command) (ResponseType, error) {
	responseTypeValue, err := c.Flags().GetString(responseTypeFlag)
	if err != nil {
		return "", err
	}

	return ResponseTypeFromString(responseTypeValue), nil
}
