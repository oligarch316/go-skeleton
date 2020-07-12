package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
)

// PrintError TODO.
type PrintError struct {
	error
	Format Format
}

// Unwrap TODO.
func (pe PrintError) Unwrap() error { return pe.error }

func (pe PrintError) Error() string {
	return fmt.Sprintf("failed to print format '%s': %s", pe.Format, pe.error)
}

// FormatPrinter TODO.
type FormatPrinter func(interface{}) ([]byte, error)

// Printer TODO.
type Printer map[string]FormatPrinter

// Marshal TODO.
func (p Printer) Marshal(format Format, v interface{}) ([]byte, error) {
	print, ok := p[format.String()]
	if !ok {
		return nil, PrintError{errors.New("unsupported format"), format}
	}

	return print(v)
}

// DefaultPrinter TODO.
func DefaultPrinter() Printer {
	return Printer{
		FormatJSON.String(): json.Marshal,
		FormatYAML.String(): yaml.Marshal,
		FormatYML.String():  yaml.Marshal,
	}
}

// Marshal TODO.
func Marshal(format Format, v interface{}) ([]byte, error) {
	return DefaultPrinter().Marshal(format, v)
}
