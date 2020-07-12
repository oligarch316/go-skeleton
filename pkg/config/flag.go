package config

import "strings"

// FileFlag TODO.
type FileFlag struct {
	Format  Format
	pattern string
}

// Readers TODO.
func (ff FileFlag) Readers() ([]Reader, error) { return singleFileFor(ff.Format, ff.pattern) }

// Set TODO.
func (ff *FileFlag) Set(s string) error {
	ff.pattern = s
	return nil
}

func (ff FileFlag) String() string { return ff.pattern }

// Type TODO.
func (ff FileFlag) Type() string { return "file" }

// MultiFileFlag TODO.
type MultiFileFlag struct {
	Format   Format
	patterns []string
}

// Readers TODO.
func (mff MultiFileFlag) Readers() ([]Reader, error) { return multiFileFor(mff.Format, mff.patterns) }

// Set TODO.
func (mff *MultiFileFlag) Set(s string) error {
	mff.patterns = append(mff.patterns, strings.Split(s, ",")...)
	return nil
}

func (mff MultiFileFlag) String() string { return strings.Join(mff.patterns, ",") }

// Type TODO.
func (mff MultiFileFlag) Type() string { return "fileSlice" }
