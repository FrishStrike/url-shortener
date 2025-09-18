package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url is not found")
	ErrURLExists   = errors.New("url exists")
)
