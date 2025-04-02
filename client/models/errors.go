package models

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized: please check your credentials")
	ErrServerIssue  = errors.New("server issue: please try again later")
	ErrBadRequest   = errors.New("bad request: invalid input provided")
	ErrConflict     = errors.New("conflict: data already exists with input provided")
	ErrNotFound     = errors.New("not found: no data with input provided")
)
