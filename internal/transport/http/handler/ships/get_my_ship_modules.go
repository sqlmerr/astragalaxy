package http_handler_ships

import (
	"net/http"

	"github.com/samber/lo"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
	http_utils "github.com/sqlmerr/astragalaxy/internal/transport/http/utils"
)

type GetMyShipModulesResponse struct {
	Data []string `json:"data"`
}

func (h *ShipsHTTPHandler) GetMyShipModules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := core_auth.GetAgentIDFromContext(ctx)
	shipID, err := http_utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get `id` UUID path value")
		return
	}

	modules, err := h.shipsService.GetShipModules(ctx, agentID, shipID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get ship modules")
		return
	}

	response := GetMyShipModulesResponse{Data: lo.Map(modules, func(i model.ShipModule, _ int) string {
		return string(i.Type)
	})}

	responseHandler.JSONResponse(http.StatusOK, response)
}
