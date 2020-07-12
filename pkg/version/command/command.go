package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Field TODO
var (
	FieldAppVersion = Field{
		Name:      "AppVersion",
		Label:     "App Version",
		ShortFlag: "a",
		LongFlag:  "app",
	}
	FieldGitRevision = Field{
		Name:      "GitRevision",
		Label:     "Git Revision",
		ShortFlag: "r",
		LongFlag:  "revision",
	}
	FieldGolangVersion = Field{
		Name:      "GolangVersion",
		Label:     "Golang Version",
		ShortFlag: "g",
		LongFlag:  "golang",
	}
)

// Field TODO
type Field struct {
	Name                string
	Label               string
	ShortFlag, LongFlag string
}

func (f Field) String() string { return fmt.Sprintf(`{{.%s}}`, f.Name) }

func (f Field) labeled() string { return fmt.Sprintf("%s: %s", f.Label, f.String()) }

func (f Field) usage() string { return fmt.Sprintf("print %s", strings.ToLower(f.Label)) }

// Options TODO
type Options struct {
	Out    io.Writer
	Fields []Field
}

type fieldFlag struct {
	Field
	onSet func(string) error
}

func (ff fieldFlag) add(fs *pflag.FlagSet) {
	if ff.LongFlag == "" {
		return
	}

	fs.VarPF(ff, ff.LongFlag, ff.ShortFlag, ff.usage()).NoOptDefVal = "true"
}

func (ff fieldFlag) IsBoolFlag() bool { return true }

func (ff fieldFlag) Type() string { return "bool" }

func (ff fieldFlag) String() string { return "false" }

func (ff fieldFlag) Set(s string) error {
	b, err := strconv.ParseBool(s)
	switch {
	case err != nil:
		return err
	case !b:
		return nil
	}

	return ff.onSet(ff.Field.String())
}

type templateString struct {
	str    string
	wasSet bool
}

func (ts templateString) Type() string { return "text/template" }

func (ts templateString) String() string { return "" }

func (ts *templateString) Set(s string) error {
	if ts.wasSet {
		return errors.New("multiple values given for template")
	}

	ts.wasSet = true
	ts.str = s
	return nil
}

// New TODO
func New(name string, v interface{}, opts ...Option) *cobra.Command {
	res := &cobra.Command{
		Use:   name,
		Short: "Print version information",
		Long:  "Print version information",
	}

	o := Options{
		Out: os.Stdout,
		Fields: []Field{
			FieldAppVersion,
			FieldGitRevision,
			FieldGolangVersion,
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	var (
		tmplString  templateString
		defaultSegs []string
	)

	for _, f := range o.Fields {
		defaultSegs = append(defaultSegs, f.labeled())
		fieldFlag{f, tmplString.Set}.add(res.Flags())
	}

	tmplString.str = strings.Join(defaultSegs, "\n")
	res.Flags().VarP(&tmplString, "template", "t", "print using custom template")

	res.RunE = func(_ *cobra.Command, _ []string) error { return run(o.Out, v, tmplString) }

	return res
}

func run(out io.Writer, v interface{}, tmplString templateString) error {
	tmpl, err := template.New("version").Parse(tmplString.str)
	if err != nil {
		return err
	}

	return tmpl.Execute(out, v)
}
