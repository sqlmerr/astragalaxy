package http_handler_mail

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

func (h *MailHTTPHandler) ClaimPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	id, err := http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, "Failed to get `id` UUID path value")
		return
	}

	agent := auth.GetAgentFromContext(ctx)

	cooldown, err := h.service.ClaimPackage(ctx, agent, id)
	if err != nil {
		rh.ErrorResponse(err, "Failed to claim package")
		return
	}

	rh.JSONResponse(http.StatusOK, http_dto.ColdownFromModel(cooldown))
}
