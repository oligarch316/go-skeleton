package config

import (
	"os"
	"strings"
)

// FileEnv TODO.
type FileEnv struct {
	Format Format
	Key    string

	pattern string
}

// Readers TODO.
func (fe FileEnv) Readers() ([]Reader, error) {
	pattern, ok := os.LookupEnv(fe.Key)
	if !ok {
		return nil, nil
	}

	return singleFileFor(fe.Format, pattern)
}

// MultiFileEnv TODO.
type MultiFileEnv struct {
	Format Format
	Key    string

	patterns []string
}

// Readers TODO.
func (mfe MultiFileEnv) Readers() ([]Reader, error) {
	raw, ok := os.LookupEnv(mfe.Key)
	if !ok {
		return nil, nil
	}

	return multiFileFor(mfe.Format, strings.Split(raw, ","))
}
