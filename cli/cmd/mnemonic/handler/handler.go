package handler

import (
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Mnemonic contains handler functions of the CLI commands related to key-vault mnemonic.
type Mnemonic struct {
	printer printer.Printer
	network core.Network
}

// New is the constructor of Mnemonic handler.
func New(printer printer.Printer, network core.Network) *Mnemonic {
	return &Mnemonic{
		printer: printer,
		network: network,
	}
}
