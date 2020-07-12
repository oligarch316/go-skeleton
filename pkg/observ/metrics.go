package observ

import (
	"fmt"
	"time"

	"github.com/oligarch316/go-skeleton/pkg/config/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/alexcesaro/statsd.v2"
)

// MetricsTagFormat TODO.
type MetricsTagFormat string

// MetricsTagFormat TODO.
const (
	MetricsTagFormatDatadog MetricsTagFormat = "datadog"
	MetricsTagFormatInflux  MetricsTagFormat = "influx"
)

func (mtf MetricsTagFormat) asOpt() (statsd.Option, error) {
	switch mtf {
	case MetricsTagFormatDatadog:
		return statsd.TagsFormat(statsd.Datadog), nil
	case MetricsTagFormatInflux:
		return statsd.TagsFormat(statsd.InfluxDB), nil
	}

	return nil, fmt.Errorf("unsupported metrics tag format '%s'", mtf)
}

// MetricsConfig TODO.
type MetricsConfig struct {
	Disabled bool `json:"disabled"`

	Address       string           `json:"address"`
	FlushPeriod   ctype.Duration   `json:"flushPeriod"`
	MaxPacketSize int              `json:"maxPacketSize"`
	GlobalPrefix  string           `json:"globalPrefix"`
	GlobalTags    []string         `json:"globalTags"`
	SampleRate    float32          `json:"sampleRate"`
	TagFormat     MetricsTagFormat `json:"tagFormat"`
}

func defaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Disabled: false,

		Address:       "localhost:8125",
		FlushPeriod:   ctype.Duration{Duration: 100 * time.Millisecond},
		MaxPacketSize: 1440,
		GlobalPrefix:  "",
		GlobalTags:    nil,
		SampleRate:    1,
		TagFormat:     MetricsTagFormatInflux,
	}
}

// Build TODO.
func (mc MetricsConfig) Build(logger *zap.Logger) (*statsd.Client, error) {
	if mc.Disabled {
		return statsd.New(statsd.Mute(true))
	}

	formatOpt, err := mc.TagFormat.asOpt()
	if err != nil {
		return nil, err
	}

	opts := []statsd.Option{
		statsd.Address(mc.Address),
		statsd.ErrorHandler(func(err error) { logger.Error("statsd error", zap.Error(err)) }),
		statsd.FlushPeriod(mc.FlushPeriod.Duration),
		statsd.MaxPacketSize(mc.MaxPacketSize),
		statsd.SampleRate(mc.SampleRate),
		formatOpt,
	}

	if mc.GlobalPrefix != "" {
		opts = append(opts, statsd.Prefix(mc.GlobalPrefix))
	}

	if len(mc.GlobalTags) > 0 {
		opts = append(opts, statsd.Tags(mc.GlobalTags...))
	}

	return statsd.New(opts...)
}

// MarshalLogObject TODO.
func (mc MetricsConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if mc.Disabled {
		enc.AddBool("disabled", true)
		return nil
	}

	enc.AddString("address", mc.Address)
	enc.AddDuration("flushPeriod", mc.FlushPeriod.Duration)
	enc.AddInt("maxPacketSize", mc.MaxPacketSize)
	enc.AddFloat32("sampleRate", mc.SampleRate)
	enc.AddString("tagFormat", string(mc.TagFormat))

	if mc.GlobalPrefix != "" {
		enc.AddString("globalPrefix", mc.GlobalPrefix)
	}

	if len(mc.GlobalTags) > 0 {
		enc.AddArray("globalTags", zapcore.ArrayMarshalerFunc(
			func(arrEnc zapcore.ArrayEncoder) error {
				for _, tag := range mc.GlobalTags {
					arrEnc.AppendString(tag)
				}
				return nil
			},
		))
	}

	return nil
}
