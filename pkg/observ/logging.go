package observ

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggingConfig TODO.
type LoggingConfig struct{ zap.Config }

func defaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Config: zap.Config{
			Development:       true,
			DisableStacktrace: true,
			Encoding:          "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			ErrorOutputPaths: []string{"stderr"},
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			OutputPaths:      []string{"stderr"},
		},
	}
}

// MarshalJSON TODO.
func (lc LoggingConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		DisableCaller     bool   `json:"disableCaller"`
		DisableStacktrace bool   `json:"disableStacktrace"`
		Encoding          string `json:"encoding"`
		Level             string `json:"level"`
	}{
		DisableCaller:     lc.DisableCaller,
		DisableStacktrace: lc.DisableStacktrace,
		Encoding:          lc.Encoding,
		Level:             lc.Level.String(),
	})
}

// MarshalLogObject TODO.
func (lc LoggingConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("disableCaller", lc.DisableCaller)
	enc.AddBool("disableStacktrace", lc.DisableStacktrace)
	enc.AddString("encoding", lc.Encoding)
	enc.AddString("level", lc.Level.String())
	return nil
}
