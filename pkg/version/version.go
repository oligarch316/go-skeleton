package version

import (
	"runtime"

	"github.com/oligarch316/go-skeleton/pkg/version/command"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

// Info contains version information.
type Info struct{ AppVersion, GitRevision, GolangVersion string }

// Init TODO.
func (i *Info) Init(appVersion, gitRevision string) {
	if appVersion == "" {
		appVersion = "unknown"
	}

	if gitRevision == "" {
		gitRevision = "unknown"
	}

	i.AppVersion = appVersion
	i.GitRevision = gitRevision
	i.GolangVersion = runtime.Version()
}

// NewCommand TODO.
func (i Info) NewCommand(name string) *cobra.Command { return command.New(name, i) }

// MarshalLogObject TODO.
func (i Info) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("appVersion", i.AppVersion)
	enc.AddString("gitRevision", i.GitRevision)
	enc.AddString("golangVersion", i.GolangVersion)
	return nil
}
