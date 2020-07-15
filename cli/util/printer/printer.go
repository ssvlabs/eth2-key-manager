package printer

// Printer represents the behavior of the CLI results printer.
type Printer interface {
	// Text prints the given text as row.
	Text(text string)

	// JSON marshales the given object into JSON string and prints that.
	JSON(obj interface{}) error

	// Error prints the given error.
	Error(err error)
}
