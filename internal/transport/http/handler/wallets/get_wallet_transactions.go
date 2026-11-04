package http_handler_wallets

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

type GetWalletTransactionsResponse struct {
	Data []TransactionDTO `json:"data"`
}

func (h *WalletsHTTPHandler) GetWalletTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	addr, err := http_utils.GetStringPathValue(r, "addr")
	if err != nil {
		rh.ErrorResponse(err, "Failed to get `addr` string path value")
		return
	}

	agentID := auth.GetAgentIDFromContext(ctx)

	transactions, err := h.service.GetPublicWalletTransactions(ctx, agentID, addr)
	if err != nil {
		rh.ErrorResponse(err, "Failed to get wallet")
		return
	}

	rh.JSONResponse(http.StatusOK, GetWalletTransactionsResponse{
		Data: lo.Map(transactions, func(i model.WalletTransaction, _ int) TransactionDTO {
			return convertTransactionModelToDTO(i)
		}),
	})
}
