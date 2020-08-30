package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth-key-manager/cli/util/cliflag"
)

// Flag names.
const (
	nameFlag    = "name"
	seedFlag    = "seed"
	storageFlag = "storage"
)

// AddNameFlag adds the name flag to the command
func AddNameFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, nameFlag, "", "account name", false)
}

// GetNameFlagValue gets the name flag from the command
func GetNameFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(nameFlag)
}

// AddSeedFlag adds the seed flag to the command
func AddSeedFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, seedFlag, "", "key-vault seed", true)
}

// GetSeedFlagValue gets the seed flag from the command
func GetSeedFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(seedFlag)
}

// AddStorageFlag adds the storage flag to the command
func AddStorageFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, storageFlag, "", "key-vault storage", true)
}

// GetStorageFlagValue gets the storage flag from the command
func GetStorageFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(storageFlag)
}
