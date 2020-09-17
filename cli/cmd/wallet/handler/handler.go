package handler

import (
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Wallet contains handler functions of the CLI commands related to wallet.
type Wallet struct {
	printer printer.Printer
	network core.Network
}

// New is the constructor of Wallet handler.
func New(printer printer.Printer, network core.Network) *Wallet {
	return &Wallet{
		printer: printer,
		network: network,
	}
}
