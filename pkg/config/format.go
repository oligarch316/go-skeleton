package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Format TODO.
type Format fmt.Stringer

// FormatString TODO.
type FormatString string

// Set TODO.
func (fs *FormatString) Set(s string) error {
	*fs = FormatString(s)
	return nil
}

func (fs FormatString) String() string { return string(fs) }

// Type TODO.
func (fs FormatString) Type() string { return "format" }

// Format TODO.
const (
	formatUnknown = ""

	FormatUnknown FormatString = formatUnknown
	FormatJSON    FormatString = "json"
	FormatYAML    FormatString = "yaml"
	FormatYML     FormatString = "yml"
)

func parseExtension(path string) Format {
	ext := filepath.Ext(path)
	if strings.HasPrefix(ext, ".") {
		return FormatString(ext[1:])
	}
	return FormatUnknown
}
