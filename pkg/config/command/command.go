package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/oligarch316/go-skeleton/pkg/config"
	"github.com/spf13/cobra"
)

// Params TODO.
type Params struct {
	Loader  config.Loader
	Printer config.Printer

	OutputFormat config.Format
	Source       config.Source

	data     interface{}
	listFlag bool
}

// New TODO.
func New(name string, v interface{}, opts ...Option) *cobra.Command {
	params := Params{
		Loader:  config.DefaultLoader(),
		Printer: config.DefaultPrinter(),

		OutputFormat: config.FormatYAML,
		Source:       config.None,

		data: v,
	}

	for _, opt := range opts {
		opt(&params)
	}

	var (
		listCmd = &cobra.Command{
			Use:   "list",
			Short: "Print list of configuration sources",
			Long:  "Print list of configuration sources",
			RunE:  func(_ *cobra.Command, _ []string) error { return runList(params.Source) },
		}

		configCmd = &cobra.Command{
			Use:   name,
			Short: "Print loaded configuration",
			Long:  "Print loaded configuration",
			RunE:  func(_ *cobra.Command, _ []string) error { return run(params) },
		}
	)

	configCmd.Flags().BoolVarP(&params.listFlag, "list", "l", false, "include configuration source list")
	configCmd.AddCommand(listCmd)

	return configCmd
}

func run(params Params) error {
	src := config.Recorder{Source: params.Source}

	if err := params.Loader.Unmarshal(&src, params.data); err != nil {
		return err
	}

	bytes, err := params.Printer.Marshal(params.OutputFormat, params.data)
	if err != nil {
		return err
	}

	if params.listFlag {
		fmt.Fprintf(os.Stdout, "Sources:\n\t%s\n\n", strings.Join(src.RecordOrNone(), "\n\t"))
	}

	_, err = os.Stdout.Write(bytes)
	return err
}

func runList(source config.Source) error {
	src := config.Recorder{Source: source}

	if _, err := src.Readers(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Sources:\n\t%s", strings.Join(src.RecordOrNone(), "\n\t"))
	return nil
}
