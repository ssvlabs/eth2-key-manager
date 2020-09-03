package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
)

// Flag names.
const (
	indexFlag   = "index"
	seedFlag    = "seed"
	storageFlag = "storage"
)

// AddIndexFlag adds the index flag to the command
func AddIndexFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, indexFlag, 0, "account index", false)
}

// GetIndexFlagValue gets the index flag from the command
func GetIndexFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(indexFlag)
}

// GetIndexFlagName gets indexFlag name
func GetIndexFlagName() string {
	return indexFlag
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
