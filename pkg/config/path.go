package config

import (
	"os"
	"path/filepath"
)

// FilePath TODO.
type FilePath struct {
	Format Format
	Path   string
}

// Readers TODO.
func (fp FilePath) Readers() ([]Reader, error) {
	file, err := os.Open(fp.Path)
	if err != nil {
		return nil, err
	}

	return []Reader{fileReader{file, fp.Format}}, nil
}

// FileGlob TODO.
type FileGlob struct {
	Format  Format
	Pattern string
}

// Matches TODO.
func (fg FileGlob) Matches() ([]Source, error) {
	paths, err := filepath.Glob(fg.Pattern)
	if err != nil {
		return nil, err
	}

	var res []Source
	for _, path := range paths {
		res = append(res, FilePath{fg.Format, path})
	}
	return res, nil
}

// Readers TODO.
func (fg FileGlob) Readers() ([]Reader, error) {
	matches, err := fg.Matches()
	if err != nil {
		return nil, err
	}
	return AllOf(matches).Readers()
}

// ----- Glob helpers

func singleFileFor(format Format, pattern string) ([]Reader, error) {
	matches, err := FileGlob{format, pattern}.Matches()
	if err != nil {
		return nil, err
	}
	return OneOf(matches).Readers()
}

func multiFileFor(format Format, patterns []string) ([]Reader, error) {
	var res []Reader
	for _, pattern := range patterns {
		readers, err := FileGlob{format, pattern}.Readers()
		if err != nil {
			return nil, err
		}
		res = append(res, readers...)
	}
	return res, nil
}
