package http_handler_wallets

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GetMyWalletsResponse struct {
	Data []WalletResponseDTO `json:"data"`
}

func (h *WalletsHTTPHandler) GetMyWallets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	agentID := auth.GetAgentIDFromContext(ctx)
	wallets, err := h.service.GetAgentWallets(ctx, agentID)
	if err != nil {
		rh.ErrorResponse(err, "Failed to get agent wallets")
		return
	}
	response := GetMyWalletsResponse{
		Data: lo.Map(wallets, func(i model.Wallet, _ int) WalletResponseDTO {
			return convertModelToDTO(i)
		}),
	}
	rh.JSONResponse(http.StatusOK, response)
}
