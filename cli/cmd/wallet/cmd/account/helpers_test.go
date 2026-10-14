package account_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// resetFlags restores every flag of c and its subcommands to its default and
// marks it unset. Subtests share cmd.RootCmd, whose flags otherwise carry over
// from one run to the next. Set(DefValue) is a true reset only for scalar flags,
// so any other kind (e.g. slices, which append) fails the test.
func resetFlags(t *testing.T, c *cobra.Command) {
	t.Helper()
	reset := func(f *pflag.Flag) {
		require.NoError(t, f.Value.Set(f.DefValue), "can't reset --%s", f.Name)
		require.Equal(t, f.DefValue, f.Value.String(), "can't reset --%s", f.Name)
		f.Changed = false
	}
	c.Flags().VisitAll(reset)
	c.PersistentFlags().VisitAll(reset)
	for _, sub := range c.Commands() {
		resetFlags(t, sub)
	}
}
