package common

import "errors"

var (
	ErrParseJSON = errors.New("failed to parse JSON")
	ErrDatabase  = errors.New("database error")
)
