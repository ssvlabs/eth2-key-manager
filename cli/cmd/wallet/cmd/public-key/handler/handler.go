package handler

import "github.com/bloxapp/eth-key-manager/cli/util/printer"

// PublicKey contains handler functions of the CLI commands related to wallet public keys.
type PublicKey struct {
	printer printer.Printer
}

// New is the constructor of PublicKey handler.
func New(printer printer.Printer) *PublicKey {
	return &PublicKey{
		printer: printer,
	}
}
