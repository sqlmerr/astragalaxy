package http_handler_wallets

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/game/wallets"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

type PayRequest struct {
	To     string `json:"to" validate:"required"`
	Units  int    `json:"units" validate:"min=0"`
	Quarks int    `json:"quarks" validate:"min=0"`
}

type PayResponse struct {
	Fee      int                  `json:"fee"`
	Cooldown http_dto.CooldownDTO `json:"cooldown"`
}

func (h *WalletsHTTPHandler) Pay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	rh := http_response.NewHTTPResponseHandler(log, w)

	var req PayRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		rh.ErrorResponse(err, "Failed to decode and validate http request json body")
		return
	}
	agentID := auth.GetAgentIDFromContext(ctx)
	fromWalletAddress, err := http_utils.GetStringPathValue(r, "addr")
	if err != nil {
		rh.ErrorResponse(err, "Failed to get `addr` string path value")
		return
	}

	output, err := h.service.Pay(ctx, wallets.PayInput{
		AgentID:     agentID,
		FromAddress: fromWalletAddress,
		ToAddress:   req.To,
		Units:       req.Units,
		Quarks:      req.Quarks,
	})
	if err != nil {
		rh.ErrorResponse(err, "Failed to pay")
		return
	}

	response := PayResponse{
		Fee:      output.Fee,
		Cooldown: http_dto.ColdownFromModel(output.Cooldown),
	}
	rh.JSONResponse(http.StatusOK, response)
}
