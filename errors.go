package crossref

import "errors"

var (
	// ErrEmptyQuery returns an error if an empty query was requested.
	ErrEmptyQuery = errors.New("empty query requested")
	// ErrZeroWorks returns an error if the query did not find any results.
	ErrZeroWorks = errors.New("no works were found")
)
