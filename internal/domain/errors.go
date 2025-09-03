package domain

import "errors"

var (
	// ErrBookNotFound is the domain error when a book is not found.
	ErrBookNotFound = errors.New("book not found")
	// ErrInvalidBookID is the domain error returned if an invalid UUID is passed.
	ErrInvalidBookID = errors.New("invalid book id")
)
