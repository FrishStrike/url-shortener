package storage

import "errors"

var (
	ErrURLExists    = errors.New("URL is exists")
	ErrURLIsNotFund = errors.New("URL is not found")
)
