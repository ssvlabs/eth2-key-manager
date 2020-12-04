package handler

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
)

// Create creates a new wallet and prints the storage
func (h *Wallet) Create(cmd *cobra.Command, _ []string) error {
	store := inmemory.NewInMemStore(h.network)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	_, err := eth2keymanager.NewKeyVault(options)
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
