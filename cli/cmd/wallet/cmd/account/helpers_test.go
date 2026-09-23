package account_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// resetFlags restores every flag of c and its subcommands to its default and
// marks it unset. Subtests share cmd.RootCmd, whose flags otherwise carry over
// from one run to the next.
func resetFlags(t *testing.T, c *cobra.Command) {
	t.Helper()
	reset := func(f *pflag.Flag) {
		require.NoError(t, f.Value.Set(f.DefValue))
		f.Changed = false
	}
	c.Flags().VisitAll(reset)
	c.PersistentFlags().VisitAll(reset)
	for _, sub := range c.Commands() {
		resetFlags(t, sub)
	}
}
