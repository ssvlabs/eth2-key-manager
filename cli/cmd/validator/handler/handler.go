package handler

import "github.com/bloxapp/eth2-key-manager/cli/util/printer"

// Handler represents the commands handler logic.
type Handler struct {
	printer printer.Printer
}

// New is the constructor of Handler.
func New(printer printer.Printer) *Handler {
	return &Handler{
		printer: printer,
	}
}
