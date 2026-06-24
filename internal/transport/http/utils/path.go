package http_utils

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, core_errors.NewWithCode(core_errors.CodeDecodeError, fmt.Errorf("no key %s in path values: %w", key, core_errors.ErrInvalidArgument))
	}

	value, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, core_errors.NewWithCode(core_errors.CodeDecodeError, fmt.Errorf(
			"path value %s by key %s not a valid UUID: %w",
			pathValue,
			key,
			core_errors.ErrInvalidArgument,
		))
	}
	return value, nil
}
