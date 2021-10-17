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
		if err := c.MarkPersistentFlagRequired(flag); err != nil {
			fmt.Printf("failed to mark persistent string flage : %s", err)
		}
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
		if err := c.MarkPersistentFlagRequired(flag); err != nil {
			fmt.Printf("failed to mark persistent int flage : %s", err)
		}
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
		if err := c.MarkPersistentFlagRequired(flag); err != nil {
			fmt.Printf("failed to mark persistent bool flage : %s", err)
		}
	}
}
