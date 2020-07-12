package command

import "github.com/oligarch316/go-skeleton/pkg/config"

// Option TODO.
type Option func(*Params)

// WithLoader TODO.
func WithLoader(loader config.Loader) Option {
	return func(p *Params) { p.Loader = loader }
}

// WithPrinter TODO.
func WithPrinter(printer config.Printer) Option {
	return func(p *Params) { p.Printer = printer }
}

// WithOutputFormat TODO.
func WithOutputFormat(format config.Format) Option {
	return func(p *Params) { p.OutputFormat = format }
}

// WithSource TODO.
func WithSource(source config.Source) Option {
	return func(p *Params) { p.Source = source }
}
