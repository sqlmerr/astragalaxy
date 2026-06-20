package http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return core_errors.NewWithCode(core_errors.CodeDecodeError, fmt.Errorf("decode json: %w: %w", core_errors.ErrInvalidArgument, err))
	}

	var err error
	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return core_errors.NewWithCode(core_errors.CodeValidationError, fmt.Errorf("request validation: %w: %w", core_errors.ErrInvalidArgument, err))
	}

	return nil
}
