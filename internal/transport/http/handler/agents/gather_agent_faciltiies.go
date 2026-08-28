package http_handler_agents

import (
	"net/http"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type GatherAgentFacilitiesResponse struct {
	Data map[string][]string `json:"data"`
}

func (h *AgentsHTTPHandler) GetAgentFacilities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	agentID := auth.GetAgentIDFromContext(ctx)
	facilities, err := h.facilitiesService.GatherAvailableFacilities(ctx, agentID)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get available facilities")
		return
	}

	response := GatherAgentFacilitiesResponse{
		Data: lo.MapKeys(facilities, func(value []string, key registry.FacilityType) string {
			return string(key)
		}),
	}
	responseHandler.JSONResponse(http.StatusOK, response)
}
