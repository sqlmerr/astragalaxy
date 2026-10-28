package http_handler_navigation

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type NavigateWarpRequest struct {
	X int `json:"x" validate:"required"`
	Y int `json:"y" validate:"required"`
}

type ResourceDTO struct {
	InventoryID  uuid.UUID `json:"inventory_id"`
	ResourceType string    `json:"resource_type"`
	Amount       int       `json:"amount"`
}

type FuelDTO struct {
	ResourceType string `json:"resource_type"`
	Used         int    `json:"used"`
	Left         int    `json:"left"`
}

type NavigateWarpResponse struct {
	Cooldown  http_dto.CooldownDTO `json:"cooldown"`
	Fuel      FuelDTO              `json:"fuel"`
	Discovery *SystemDiscoveryDTO  `json:"discovery"`
}

func (h *NavigationHTTPHandler) NavigateWarp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var req NavigateWarpRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request body")
		return
	}

	agentID := auth.GetAgentIDFromContext(ctx)

	cooldown, fuelUsage, discovery, err := h.navigationService.NavigateWarp(ctx, agentID, req.X, req.Y)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to warp")
		return
	}

	var discoveryDTO *SystemDiscoveryDTO
	if discovery != nil {
		discoveryDTO = &SystemDiscoveryDTO{
			SystemX:      discovery.SystemX,
			SystemY:      discovery.SystemY,
			DiscoveredAt: discovery.DiscoveredAt,
			DiscoveredBy: discovery.DiscoveredBy,
		}
	}

	response := NavigateWarpResponse{
		Cooldown: http_dto.ColdownFromModel(cooldown),
		Fuel: FuelDTO{
			ResourceType: string(fuelUsage.ResourceType),
			Used:         fuelUsage.Used,
			Left:         fuelUsage.Left,
		},
		Discovery: discoveryDTO,
	}
	responseHandler.JSONResponse(http.StatusOK, response)
}
