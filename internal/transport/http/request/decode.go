package http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return errs.NewWithCode(errs.CodeDecodeError, fmt.Errorf("decode json: %w: %w", errs.ErrInvalidArgument, err))
	}

	var err error
	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return errs.NewWithCode(errs.CodeValidationError, fmt.Errorf("request validation: %w: %w", errs.ErrInvalidArgument, err))
	}

	return nil
}
