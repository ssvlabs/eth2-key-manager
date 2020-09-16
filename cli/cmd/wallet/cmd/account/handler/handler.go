package handler

import (
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
	"github.com/bloxapp/eth2-key-manager/core"
)

// Account contains handler functions of the CLI commands related to wallet account.
type Account struct {
	printer printer.Printer
	network core.Network
}

// New is the constructor of Account handler.
func New(printer printer.Printer, network core.Network) *Account {
	return &Account{
		printer: printer,
		network: network,
	}
}
