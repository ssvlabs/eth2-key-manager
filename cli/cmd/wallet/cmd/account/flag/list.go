package flag

import (
	"github.com/spf13/cobra"

	"github.com/ssvlabs/eth2-key-manager/cli/util/cliflag"
)

// Flag names.
const (
	storageFlag = "storage"
)

// AddStorageFlag adds the storage flag to the command
func AddStorageFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, storageFlag, "", "key-vault storage", true)
}

// GetStorageFlagValue gets the storage flag from the command
func GetStorageFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(storageFlag)
}
