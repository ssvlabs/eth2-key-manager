package handler

import "github.com/bloxapp/eth2-key-manager/cli/util/printer"

// Wallet contains handler functions of the CLI commands related to wallet.
type Wallet struct {
	printer printer.Printer
}

// New is the constructor of Wallet handler.
func New(printer printer.Printer) *Wallet {
	return &Wallet{
		printer: printer,
	}
}
