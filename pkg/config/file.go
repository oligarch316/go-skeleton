package config

import "os"

// ===== Standard input

const stdinName = "standard input"

type stdinReader struct{ format Format }

func (sr stdinReader) Format() string { return sr.format.String() }

func (stdinReader) Name() string { return stdinName }

func (sr stdinReader) Read(p []byte) (int, error) { return os.Stdin.Read(p) }

// Stdin TODO.
type Stdin struct{ Format Format }

// Readers TODO.
func (s Stdin) Readers() ([]Reader, error) { return []Reader{stdinReader{s.Format}}, nil }

// ===== General file

type fileReader struct {
	*os.File
	format Format
}

func (fr fileReader) Format() string {
	if fr.format != nil {
		if res := fr.format.String(); res != formatUnknown {
			return res
		}
	}

	return parseExtension(fr.Name()).String()
}

// File TODO.
type File struct {
	*os.File
	Format Format
}

// Readers TODO.
func (f File) Readers() ([]Reader, error) { return []Reader{fileReader{f.File, f.Format}}, nil }
