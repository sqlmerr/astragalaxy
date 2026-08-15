package http_utils

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, errs.NewWithCode(errs.CodeDecodeError, fmt.Errorf("no key %s in path values: %w", key, errs.ErrInvalidArgument))
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
