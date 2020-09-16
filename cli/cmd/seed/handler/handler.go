package handler

import (
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Seed contains handler functions of the CLI commands related to key-vault seed.
type Seed struct {
	printer printer.Printer
	network core.Network
}

// New is the constructor of Seed handler.
func New(printer printer.Printer, network core.Network) *Seed {
	return &Seed{
		printer: printer,
		network: network,
	}
}
