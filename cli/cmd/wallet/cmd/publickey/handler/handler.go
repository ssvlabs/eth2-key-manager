package handler

import (
	"github.com/ssvlabs/eth2-key-manager/cli/util/printer"
	"github.com/ssvlabs/eth2-key-manager/core"
)

// PublicKey contains handler functions of the CLI commands related to wallet public keys.
type PublicKey struct {
	printer printer.Printer
	network core.Network
}

// New is the constructor of PublicKey handler.
func New(printer printer.Printer, network core.Network) *PublicKey {
	return &PublicKey{
		printer: printer,
		network: network,
	}
}
