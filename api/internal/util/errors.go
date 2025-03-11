package util

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
	ErrInvalidToken                   = New("invalid token", http.StatusForbidden)
	ErrUnauthorized                   = New("unauthorized", http.StatusUnauthorized)
	ErrNotFound                       = New("not found", http.StatusNotFound)
	ErrSpaceshipNotFound              = ErrNotFound
	ErrUserNotFound                   = ErrNotFound
	ErrPlanetNotFound                 = ErrNotFound
	ErrItemNotFound                   = ErrNotFound
	ErrItemDataTagNotFound            = ErrNotFound
	ErrIDMustBeUUID                   = New("id must be an uuid type", 400)
	ErrSpaceshipAlreadyFlying         = New("spaceship is already flying", http.StatusBadRequest)
	ErrSpaceshipIsInAnotherSystem     = New("spaceship is in another system", http.StatusBadRequest)
	ErrSpaceshipIsAlreadyInThisPlanet = New("spaceship is already in this planet", http.StatusBadRequest)
	ErrSpaceshipIsAlreadyInThisSystem = New("spaceship is already in this system", http.StatusBadRequest)
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
