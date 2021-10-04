package handler

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/core"
)

// Create creates a new wallet account(s) and prints the storage.
func (h *Account) Create(cmd *cobra.Command, args []string) error {
	err := core.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	accountFlags, err := CollectAccountFlags(cmd)
	if err != nil {
		return err
	}

	err = h.BuildAccounts(accountFlags)
	if err != nil {
		return err
	}
	return nil
}
