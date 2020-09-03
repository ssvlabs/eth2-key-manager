package handler

import "github.com/bloxapp/eth2-key-manager/cli/util/printer"

// Seed contains handler functions of the CLI commands related to key-vault seed.
type Seed struct {
	printer printer.Printer
}

// New is the constructor of Seed handler.
func New(printer printer.Printer) *Seed {
	return &Seed{
		printer: printer,
	}
}
