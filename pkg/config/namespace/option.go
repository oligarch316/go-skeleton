package namespace

import (
	"github.com/oligarch316/go-skeleton/pkg/config"
)

// Option TODO.
type Option func(*Params)

// Params TODO.
type Params struct {
	DefaultSourceFormat config.FormatString
	DefaultTargetFormat config.FormatString
	XDGFilePaths        []string

	Loader  config.Loader
	Printer config.Printer
}

// WithDefaultSourceFormat TODO.
func WithDefaultSourceFormat(format config.FormatString) Option {
	return func(p *Params) { p.DefaultSourceFormat = format }
}

// WithDefaultTargetFormat TODO.
func WithDefaultTargetFormat(format config.FormatString) Option {
	return func(p *Params) { p.DefaultTargetFormat = format }
}

// WithXDGFilePaths TODO.
func WithXDGFilePaths(paths ...string) Option {
	return func(p *Params) { p.XDGFilePaths = paths }
}

// WithLoader TODO.
func WithLoader(loader config.Loader) Option {
	return func(p *Params) { p.Loader = loader }
}

// WithPrinter TODO.
func WithPrinter(printer config.Printer) Option {
	return func(p *Params) { p.Printer = printer }
}
