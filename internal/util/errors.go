package util

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status_code"`
}

// Error implements error.
func (e APIError) Error() string {
	return e.Message
}

func (e APIError) GetStatus() int {
	return e.Status
}

func New(msg string, status int) APIError {
	return APIError{Message: msg, Status: status}
}

var (
	ErrUserAlreadyExists              = New("user with this username already exists", http.StatusConflict)
	ErrAstralAlreadyExists            = New("astral with this code already exists", http.StatusConflict)
	ErrServerError                    = New("server error", http.StatusInternalServerError)
	ErrInvalidToken                   = New("invalid token", http.StatusForbidden)
	ErrUnauthorized                   = New("unauthorized", http.StatusUnauthorized)
	ErrNotFound                       = New("not found", http.StatusNotFound)
	ErrSpaceshipNotFound              = ErrNotFound
	ErrPlanetNotFound                 = ErrNotFound
	ErrItemNotFound                   = ErrNotFound
	ErrItemDataTagNotFound            = ErrNotFound
	ErrIDMustBeUUID                   = New("id must be valid uuid", 400)
	ErrSpaceshipAlreadyFlying         = New("spaceship is already flying", http.StatusBadRequest)
	ErrSpaceshipIsFlying              = New("spaceship is flying", 400)
	ErrSpaceshipIsInAnotherSystem     = New("spaceship is in another system", http.StatusBadRequest)
	ErrSpaceshipIsAlreadyInThisPlanet = New("spaceship is already in this planet", http.StatusBadRequest)
	ErrSpaceshipIsAlreadyInThisSystem = New("spaceship is already in this system", http.StatusBadRequest)
	ErrPlayerAlreadyInSpaceship       = New("player already in spaceship", http.StatusBadRequest)
	ErrPlayerNotInSpaceship           = New("player not in spaceship", http.StatusBadRequest)
	ErrPlayerMustBeOutOfSpaceship     = New("player must be out of spaceship", 400)
	ErrInvalidHyperJumpPath           = New("invalid hyperjump path", http.StatusBadRequest)
	ErrInvalidCode                    = New("invalid code", http.StatusBadRequest)
	ErrTooManyAstrals                 = New("too many astrals", http.StatusBadRequest)
	ErrInvalidAstralIDHeader          = New("astral id header not specified or invalid", http.StatusUnauthorized)
	ErrAlreadyExploring               = New("you are already exploring", 400)
	ErrInvalidExplorationType = New("invalid exploration type", 400)
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
			Message: apiErr.Message,
			Status:  apiErr.Status,
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

func WriteError(api huma.API, ctx huma.Context, err error) {
	var apiErr APIError
	ok := errors.As(err, &apiErr)
	if ok {
		huma.WriteErr(api, ctx, apiErr.Status, apiErr.Message, err)
	} else {
		huma.WriteErr(api, ctx, 500, err.Error(), err)
	}
}
