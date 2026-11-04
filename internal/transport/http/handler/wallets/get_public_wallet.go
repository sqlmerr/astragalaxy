package http_handler_wallets

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

func (h *WalletsHTTPHandler) GetPublicWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	addr, err := http_utils.GetStringPathValue(r, "addr")
	if err != nil {
		rh.ErrorResponse(err, "Failed to get `addr` string path value")
		return
	}

	agentID := auth.GetAgentIDFromContext(ctx)

	wallet, err := h.service.GetPublicWallet(ctx, agentID, addr)
	if err != nil {
		rh.ErrorResponse(err, "Failed to get wallet")
		return
	}

	rh.JSONResponse(http.StatusOK, convertWalletModelToDTO(wallet))
}
