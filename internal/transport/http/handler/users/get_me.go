package http_handler_users

import (
	"net/http"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func (h *UsersHTTPHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userID := core_auth.GetUserIDFromContext(ctx)
	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get user")
		return
	}

	dto := userDTOFromModel(user)
	responseHandler.JSONResponse(http.StatusOK, dto)
}
