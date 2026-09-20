package http_utils

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

func GetStringPathValue(r *http.Request, key string) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", errs.NewWithCode(errs.CodeDecodeError, fmt.Errorf("no key %s in path values: %w", key, errs.ErrInvalidArgument))
	}

	return pathValue, nil
}

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue, err := GetStringPathValue(r, key)
	if err != nil {
		return uuid.Nil, err
	}

	value, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, errs.NewWithCode(errs.CodeDecodeError, fmt.Errorf(
			"path value %s by key %s not a valid UUID: %w",
			pathValue,
			key,
			errs.ErrInvalidArgument,
		))
	}
	return value, nil
}
