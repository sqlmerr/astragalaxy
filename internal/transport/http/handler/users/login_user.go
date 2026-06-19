package http_handler_users

import (
	"net/http"

	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (h *UsersHTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var request LoginUserRequest
	if err := http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request")
		return
	}

	accessToken, err := h.service.LoginUser(ctx, request.Username, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to login")
		return
	}

	response := LoginUserResponse{AccessToken: accessToken, TokenType: "Bearer"}
	responseHandler.JSONResponse(http.StatusOK, response)
}
