package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// StandardPrinter implements Printer interface.
// Prints stuff to the out.
type StandardPrinter struct {
	out io.Writer
}

// NewStandardOutputPrinter is the constructor of StandardOutputPrinter.
// Uses stdout to print data.
func NewStandardOutputPrinter() Printer {
	return &StandardPrinter{
		out: os.Stdout,
	}
}

// New is the constructor of StandardOutputPrinter.
// Uses the given writer to print data.
func New(out io.Writer) Printer {
	return &StandardPrinter{
		out: out,
	}
}

// Text implements Printer interface.
func (p *StandardPrinter) Text(text string) {
	fmt.Fprintln(p.out, text)
}

// JSON implements Printer interface.
func (p *StandardPrinter) JSON(obj interface{}) error {
	encoder := json.NewEncoder(p.out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(obj); err != nil {
		return errors.Wrapf(err, "failed to marshal the given object with type %T", obj)
	}

	return nil
}

// JSON implements Printer interface.
func (p *StandardPrinter) Error(err error) {
	if err != nil {
		fmt.Fprintln(p.out, "Error:", err.Error())
	}
}
