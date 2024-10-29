package handler

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ssvlabs/eth2-key-manager/cli/util/printer"
	"github.com/ssvlabs/eth2-key-manager/core"
)

// Config contains handler functions of the CLI commands related to key-vault config.
type Config struct {
	printer printer.Printer
	network core.Network
}

// New is the constructor of Seed handler.
func New(printer printer.Printer, network core.Network) *Config {
	return &Config{
		printer: printer,
		network: network,
	}
}

// CurrentSlot prints the current slot.
func (h *Config) CurrentSlot(_ *cobra.Command, _ []string) error {
	h.printer.Text(fmt.Sprintf("%d", h.network.EstimatedCurrentSlot()))
	return nil
}

// CurrentEpoch prints the current epoch.
func (h *Config) CurrentEpoch(_ *cobra.Command, _ []string) error {
	h.printer.Text(fmt.Sprintf("%d", h.network.EstimatedCurrentEpoch()))
	return nil
}
