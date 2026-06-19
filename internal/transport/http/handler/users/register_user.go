package http_handler_users

import (
	"net/http"

	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=1"`
}

func (h *UsersHTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request RegisterUserRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP response")
		return
	}

	user, err := h.service.RegisterUser(ctx, request.Username, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to register user")
		return
	}

	dto := userDTOFromModel(user)
	responseHandler.JSONResponse(http.StatusCreated, dto)
}
