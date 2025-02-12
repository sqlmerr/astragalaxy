package utils

import (
	"errors"
	"net/http"
)

type APIError struct {
	message string
	status  int
}

// Error implements error.
func (e APIError) Error() string {
	return e.message
}

func (e APIError) Status() int {
	return e.status
}

func New(msg string, status int) APIError {
	return APIError{message: msg, status: status}
}

var (
	ErrUserAlreadyExists              = New("user with this username already exists", http.StatusConflict)
	ErrServerError                    = New("server error", http.StatusInternalServerError)
	ErrInvalidToken                   = New("invalid token (for example 123456789:AbCdEf123456xyz7v)", http.StatusForbidden)
	ErrUnauthorized                   = New("unauthorized", http.StatusUnauthorized)
	ErrSpaceshipNotFound              = New("spaceship not found", http.StatusNotFound)
	ErrUserNotFound                   = New("user not found", http.StatusNotFound)
	ErrPlanetNotFound                 = New("planet not found", http.StatusNotFound)
	ErrItemNotFound                   = New("item not found", http.StatusNotFound)
	ErrItemDataTagNotFound            = New("item data tag not found", http.StatusNotFound)
	ErrSpaceshipAlreadyFlying         = New("spaceship is already flying", http.StatusBadRequest)
	ErrSpaceshipIsInAnotherSystem     = New("spaceship is in another system", http.StatusBadRequest)
	ErrSpaceshipIsAlreadyInThisPlanet = New("spaceship is already in this planet", http.StatusBadRequest)
	ErrPlayerAlreadyInSpaceship       = New("player already in spaceship", http.StatusBadRequest)
	ErrPlayerNotInSpaceship           = New("player not in spaceship", http.StatusBadRequest)
)

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status_code"`
}

func NewError(err error) Error {
	var apiErr APIError
	ok := errors.As(err, &apiErr)
	if ok {
		e := Error{
			Message: apiErr.message,
			Status:  apiErr.status,
		}
		return e
	}

	e := Error{}
	e.Status = 500
	switch v := err.(type) {
	default:
		e.Message = v.Error()
	}
	return e
}
