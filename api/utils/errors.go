package utils

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this username already exists")
	ErrServerError       = errors.New("server error")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUnauthorized      = errors.New("unauthorized")
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

// add switch other variant
func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	switch v := err.(type) {
	default:
		e.Errors["body"] = v.Error()
	}
	return e
}
