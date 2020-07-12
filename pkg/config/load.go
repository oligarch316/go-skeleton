package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// FormatLoader TODO.
type FormatLoader func([]byte, interface{}) error

// Loader TODO.
type Loader map[string]FormatLoader

// Unmarshal TODO.
func (l Loader) Unmarshal(source Source, v interface{}) error {
	if source == nil {
		return errors.New("nil source")
	}

	readers, err := source.Readers()
	if err != nil {
		return fmt.Errorf("failed to load source: %w", err)
	}

	for _, reader := range readers {
		load, ok := l[reader.Format()]
		if !ok {
			err = fmt.Errorf("%s: unsupported format '%s'", reader.Name(), reader.Format())
			break
		}

		data, dataErr := ioutil.ReadAll(reader)
		if dataErr != nil {
			err = fmt.Errorf("%s: failed to read data: %w", reader.Name(), dataErr)
			break
		}

		if err = load(data, v); err != nil {
			err = fmt.Errorf("%s: failed to load data: %w", reader.Name(), err)
			break
		}
	}

	for _, reader := range readers {
		if closer, ok := reader.(io.Closer); ok {
			if closeErr := closer.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}
	}

	return err
}

// DefaultLoader TODO.
func DefaultLoader() Loader {
	return Loader{
		FormatJSON.String(): json.Unmarshal,
		FormatYAML.String(): yaml.Unmarshal,
		FormatYML.String():  yaml.Unmarshal,
	}
}

// Unmarshal TODO.
func Unmarshal(source Source, v interface{}) error {
	return DefaultLoader().Unmarshal(source, v)
}
