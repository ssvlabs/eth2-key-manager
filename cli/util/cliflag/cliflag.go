package cliflag

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AddPersistentStringFlag adds a string flag to the command
func AddPersistentStringFlag(c *cobra.Command, flag string, value string, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().String(flag, value, fmt.Sprintf("%s%s", description, req))

	if isRequired {
		c.MarkPersistentFlagRequired(flag)
	}
}

// AddPersistentIntFlag adds a integer flag to the command
func AddPersistentIntFlag(c *cobra.Command, flag string, value int, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().Int(flag, value, fmt.Sprintf("%s%s", description, req))

	if isRequired {
		c.MarkPersistentFlagRequired(flag)
	}
}

// AddPersistentInt64SliceFlag adds a int64 slice flag to the command
func AddPersistentInt64SliceFlag(c *cobra.Command, flag string, value []int64, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().Int64Slice(flag, value, fmt.Sprintf("%s%s", description, req))

	if isRequired {
		c.MarkPersistentFlagRequired(flag)
	}
}

// AddPersistentBoolFlag adds a bool flag to the command
func AddPersistentBoolFlag(c *cobra.Command, flag string, value bool, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().Bool(flag, value, fmt.Sprintf("%s%s", description, req))

	if isRequired {
		c.MarkPersistentFlagRequired(flag)
	}
}
