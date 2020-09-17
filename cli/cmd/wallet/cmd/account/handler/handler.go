package handler

import (
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

// Account contains handler functions of the CLI commands related to wallet account.
type Account struct {
	printer printer.Printer
}

// New is the constructor of Account handler.
func New(printer printer.Printer) *Account {
	return &Account{
		printer: printer,
	}
}
