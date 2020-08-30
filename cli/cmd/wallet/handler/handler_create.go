package handler

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/bloxapp/eth-key-manager"
	"github.com/bloxapp/eth-key-manager/stores/in_memory"
)

// Wallet creates a new wallet and prints the storage
func (h *Wallet) Create(cmd *cobra.Command, args []string) error {
	store := in_memory.NewInMemStore()
	options := &KeyVault.KeyVaultOptions{}
	options.SetStorage(store)

	_, err := KeyVault.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault.")
	}

	// marshal storage
	bytes, err := store.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "failed to JSON marshal storage.")
	}

	h.printer.Text(hex.EncodeToString(bytes))
	return nil
}
