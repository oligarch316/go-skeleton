package namespace

import (
	"path/filepath"
	"strings"

	"github.com/oligarch316/go-skeleton/pkg/config"
	"github.com/oligarch316/go-skeleton/pkg/config/command"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	flagNameOutputFormat = "out-format"
	flagNameSource       = "source"
	flagNameSourceFormat = "source-format"

	usageOutputFormat = "configuration output format"
	usageSource       = "configuration source"
	usageSourceFormat = "configuration source format"
)

// NS TODO.
type NS struct {
	config.Source

	Name    string
	Loader  config.Loader
	Printer config.Printer

	defaultTargetFormat config.FormatString
	sourceFlag          *config.FileFlag
	sourceFormatFlag    *config.FormatString
}

func formatEnvVar(segs ...string) string {
	return strings.ToUpper(strings.Join(segs, "_"))
}

func formatFlagName(segs ...string) string {
	return strings.ToLower(strings.Join(segs, "-"))
}

// New TODO.
func New(name string, opts ...Option) *NS {
	params := Params{
		DefaultSourceFormat: config.FormatUnknown,
		DefaultTargetFormat: config.FormatYAML,
		XDGFilePaths:        []string{"config.yaml", "config.yml", "config.json"},

		Loader:  config.DefaultLoader(),
		Printer: config.DefaultPrinter(),
	}

	for _, opt := range opts {
		opt(&params)
	}

	format := params.DefaultSourceFormat

	fileFlag := config.FileFlag{Format: &format}
	fileEnv := config.FileEnv{Format: &format, Key: formatEnvVar(name, "config", "source")}

	xdgPaths := config.XDGPaths{}
	for _, fileName := range params.XDGFilePaths {
		newPath := config.FilePath{Path: filepath.Join(name, fileName)}
		xdgPaths.RelativePaths = append(xdgPaths.RelativePaths, newPath)
	}

	return &NS{
		Source: config.OneOf{
			&fileFlag,
			fileEnv,
			xdgPaths,
			config.None,
		},

		Name:    name,
		Loader:  params.Loader,
		Printer: params.Printer,

		defaultTargetFormat: params.DefaultTargetFormat,
		sourceFlag:          &fileFlag,
		sourceFormatFlag:    &format,
	}
}

func (ns NS) String() string { return ns.Name }

func (ns *NS) setSourceFlags(fs *pflag.FlagSet, srcName, srcFmtName string) {
	fs.Var(ns.sourceFlag, srcName, usageSource)
	fs.Var(ns.sourceFormatFlag, srcFmtName, usageSourceFormat)
}

// SetFlags TODO
func (ns *NS) SetFlags(fs *pflag.FlagSet) {
	ns.setSourceFlags(
		fs,
		formatFlagName("config", flagNameSource),
		formatFlagName("config", flagNameSourceFormat),
	)
}

// Load TODO.
func (ns *NS) Load(defaults interface{}) error {
	return ns.Loader.Unmarshal(ns.Source, defaults)
}

// LoadAndRecord TODO.
func (ns *NS) LoadAndRecord(defaults interface{}) ([]string, error) {
	src := config.Recorder{Source: ns.Source}
	err := ns.Loader.Unmarshal(&src, defaults)
	return src.Record, err
}

// NewCommand TODO.
func (ns *NS) NewCommand(name string, defaults interface{}) *cobra.Command {
	outFormat := ns.defaultTargetFormat

	res := command.New(
		name,
		defaults,
		command.WithSource(ns.Source),
		command.WithOutputFormat(&outFormat),
		command.WithLoader(ns.Loader),
		command.WithPrinter(ns.Printer),
	)

	fs := res.Flags()
	ns.setSourceFlags(fs, flagNameSource, flagNameSourceFormat)
	fs.Var(&outFormat, flagNameOutputFormat, usageOutputFormat)

	return res
}
