package flag

import (
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
)

// Flag names.
const (
	publicKeyFlag = "public-key"
)

// AddPublicKeyFlag adds the public key flag to the command
func AddPublicKeyFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, publicKeyFlag, "", "public key", true)
}

// GetPublicKeyFlagValue gets the public key flag from the command
func GetPublicKeyFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(publicKeyFlag)
}
