package http_handler_navigation

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type NavigatePlanetRequest struct {
	Orbit *int `json:"orbit" validate:"required"`
}

type NavigatePlanetResponse struct {
	NavigationResponseDTO
	Discovery *PlanetDiscoveryDTO `json:"discovery"`
}

func (h *NavigationHTTPHandler) NavigatePlanet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var req NavigatePlanetRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request body")
		return
	}

	agentID := auth.GetAgentIDFromContext(ctx)

	cooldown, discovery, err := h.navigationService.NavigatePlanet(ctx, agentID, *req.Orbit)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to navigate")
		return
	}

	var discoveryDTO *PlanetDiscoveryDTO
	if discovery != nil {
		discoveryDTO = &PlanetDiscoveryDTO{
			SystemX:      discovery.SystemX,
			SystemY:      discovery.SystemY,
			Orbit:        discovery.Orbit,
			DiscoveredAt: discovery.DiscoveredAt,
			DiscoveredBy: discovery.DiscoveredBy,
		}
	}

	response := NavigatePlanetResponse{
		NavigationResponseDTO: NavigationResponseDTO{Cooldown: http_dto.ColdownFromModel(cooldown)},
		Discovery:             discoveryDTO,
	}
	responseHandler.JSONResponse(http.StatusOK, response)
}
