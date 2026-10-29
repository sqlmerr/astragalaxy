package http_handler_wallets

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func (h *WalletsHTTPHandler) RegisterWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	agentID := auth.GetAgentIDFromContext(ctx)

	wallet, err := h.service.RegisterWallet(ctx, agentID)
	if err != nil {
		rh.ErrorResponse(err, "Failed to register a new wallet")
		return
	}

	rh.JSONResponse(http.StatusCreated, convertModelToDTO(wallet))
}
