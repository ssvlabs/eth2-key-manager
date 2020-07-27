package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/KeyVault/cli/util/cliflag"
)

// Flag names.
const (
	indexFlag = "index"
	seedFlag = "seed"
)

// AddIndexFlag adds the index flag to the command
func AddIndexFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, indexFlag, 0, "account index", true)
}

// GetIndexFlagValue gets the index flag from the command
func GetIndexFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(indexFlag)
}

// AddSeedFlag adds the seed flag to the command
func AddSeedFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, seedFlag, "", "key-vault seed", true)
}

// GetSeedFlagValue gets the seed flag from the command
func GetSeedFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(seedFlag)
}
