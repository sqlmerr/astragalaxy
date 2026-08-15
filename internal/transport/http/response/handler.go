package http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *core_logger.Logger, w http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log, w: w,
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v: %w", p, errs.ErrInternal)

	h.ErrorResponse(err, msg)
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)
	switch {
	case errors.Is(err, errs.ErrInternal):
		statusCode = http.StatusInternalServerError
		err = errs.NewWithCode(errs.CodeInternalServerError, err)
		logFunc = h.log.Error
	case errors.Is(err, errs.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, errs.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	case errors.Is(err, errs.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	case errors.Is(err, errs.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn
	case errors.Is(err, errs.ErrAccessDenied):
		statusCode = http.StatusForbidden
		logFunc = h.log.Warn
	case errors.Is(err, errs.ErrUnprocessableEntity):
		statusCode = http.StatusUnprocessableEntity
		logFunc = h.log.Warn
	case errors.Is(err, errs.ErrNotModified):
		statusCode = http.StatusNotModified
		logFunc = h.log.Debug
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))
	if statusCode == http.StatusInternalServerError {
		err = errs.NewWithCode(errs.CodeInternalServerError, errs.ErrInternal)
	}
	var errorCode errs.ErrorCode
	if errorWithCode, ok := errors.AsType[errs.WithCode](err); ok {
		errorCode = errorWithCode.Code
	} else {
		errorCode = errs.CodeUnknown
	}
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
		"code":    string(errorCode),
	}
	h.JSONResponse(statusCode, response)
}

func (h *HTTPResponseHandler) JSONResponse(statusCode int, data any) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(data); err != nil {
		h.log.Error("write HTTP Response", zap.Error(err))
	}
}
