package handler

import (
	"io"

	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
	"github.com/bloxapp/eth2-key-manager/core"
)

// ResultWriterFactory is the create validator result factory
type ResultWriterFactory func(name string) (io.Writer, func(), error)

// Handler represents the commands handler logic.
type Handler struct {
	printer             printer.Printer
	resultWriterFactory ResultWriterFactory
	network             core.Network
}

// New is the constructor of Handler.
func New(printer printer.Printer, resultWriterFactory ResultWriterFactory, network core.Network) *Handler {
	return &Handler{
		printer:             printer,
		resultWriterFactory: resultWriterFactory,
		network:             network,
	}
}
