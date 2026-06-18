package core_errors

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrInternal            = errors.New("server error")
	ErrConflict            = errors.New("conflict")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrAccessDenied        = errors.New("access denied")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
)
