package handler

import "github.com/bloxapp/eth-key-manager/cli/util/printer"

// Mnemonic contains handler functions of the CLI commands related to key-vault mnemonic.
type Mnemonic struct {
	printer printer.Printer
}

// New is the constructor of Mnemonic handler.
func New(printer printer.Printer) *Mnemonic {
	return &Mnemonic{
		printer: printer,
	}
}
