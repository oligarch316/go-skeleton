package observ

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/alexcesaro/statsd.v2"
)

// Config TODO.
type Config struct {
	Logging LoggingConfig `json:"logging"`
	Metrics MetricsConfig `json:"metrics"`
}

// DefaultConfig TODO.
func DefaultConfig() Config {
	return Config{
		Logging: defaultLoggingConfig(),
		Metrics: defaultMetricsConfig(),
	}
}

// MarshalLogObject TODO.
func (c Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("logging", c.Logging)
	enc.AddObject("metrics", c.Metrics)
	return nil
}

// Core TODO.
type Core struct{ Corelet }

// Noop TODO.
func Noop() *Core {
	logger := zap.NewNop()
	emitter, _ := statsd.New(statsd.Mute(true))

	return &Core{
		Corelet: Corelet{
			Logger:  logger,
			Emitter: emitter,
		},
	}
}

// New TODO.
func New(cfg Config) (*Core, error) {
	logger, err := cfg.Logging.Build()
	if err != nil {
		return nil, err
	}

	emitter, err := cfg.Metrics.Build(logger.Named("emitter"))
	if err != nil {
		return nil, err
	}

	return &Core{
		Corelet: Corelet{
			Logger:  logger,
			Emitter: emitter,
		},
	}, nil
}

// Bootstrap TODO.
func Bootstrap(cfg Config) (*Core, func() error, error) {
	core, err := New(cfg)
	if err != nil {
		return nil, nil, err
	}

	resetStdLog, err := zap.RedirectStdLogAt(core.Named("stdlog").Logger, zap.ErrorLevel)
	if err != nil {
		core.Close()
		return nil, nil, err
	}

	cleanup := func() error {
		resetStdLog()
		return core.Close()
	}

	return core, cleanup, nil
}

// Close TODO.
func (c *Core) Close() error {
	// FIXME: Rectify sync issues for /dev/[stdout|stderr], allowing return of .Sync() errors here
	c.Logger.Sync()
	c.Emitter.Close()
	return nil
}

// Corelet TODO.
type Corelet struct {
	Logger  *zap.Logger
	Emitter *statsd.Client
}

// Named TODO.
func (c *Corelet) Named(name string) *Corelet {
	return &Corelet{
		Logger:  c.Logger.Named(name),
		Emitter: c.Emitter.Clone(statsd.Prefix(name)),
	}
}
