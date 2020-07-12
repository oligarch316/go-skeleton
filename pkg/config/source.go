package config

import (
	"errors"
	"io"
)

type (
	// Reader TODO.
	Reader interface {
		io.Reader
		Format() string
		Name() string
	}

	// Source TODO.
	Source interface{ Readers() ([]Reader, error) }
)

// None TODO.
const None sourceNone = 0

type sourceNone int

// Readers TODO.
func (sourceNone) Readers() ([]Reader, error) { return nil, nil }

type (
	// AllOf TODO.
	AllOf []Source

	// OneOf TODO.
	OneOf []Source

	// Required TODO.
	Required struct{ Source }
)

// Readers TODO.
func (ao AllOf) Readers() ([]Reader, error) {
	var res []Reader
	for _, source := range ao {
		readers, err := source.Readers()
		if err != nil {
			return nil, err
		}
		res = append(res, readers...)
	}
	return res, nil
}

// Readers TODO.
func (oo OneOf) Readers() ([]Reader, error) {
	for _, source := range oo {
		readers, err := source.Readers()

		switch {
		case err != nil:
			return nil, err
		case readers == nil:
			continue
		}

		return readers, nil
	}

	return nil, nil
}

// Readers TODO.
func (r Required) Readers() ([]Reader, error) {
	res, err := r.Source.Readers()

	switch {
	case err != nil:
		return nil, err
	case res == nil:
		return nil, errors.New("no valid readers")
	}

	return res, nil
}

// Recorder TODO.
type Recorder struct {
	Source Source
	Record []string
}

// Readers TODO.
func (r *Recorder) Readers() ([]Reader, error) {
	res, err := r.Source.Readers()
	for _, rdr := range res {
		r.Record = append(r.Record, rdr.Name())
	}
	return res, err
}

// RecordOrNone TODO.
func (r Recorder) RecordOrNone() []string {
	if len(r.Record) < 1 {
		return []string{"None"}
	}
	return r.Record
}
