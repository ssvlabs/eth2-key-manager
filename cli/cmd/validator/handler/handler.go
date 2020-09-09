package handler

import (
	"io"

	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

// ResultWriterFactory is the create validator result factory
type ResultWriterFactory func(name string) (io.Writer, func(), error)

// Handler represents the commands handler logic.
type Handler struct {
	printer             printer.Printer
	resultWriterFactory ResultWriterFactory
}

// New is the constructor of Handler.
func New(printer printer.Printer, resultWriterFactory ResultWriterFactory) *Handler {
	return &Handler{
		printer:             printer,
		resultWriterFactory: resultWriterFactory,
	}
}
