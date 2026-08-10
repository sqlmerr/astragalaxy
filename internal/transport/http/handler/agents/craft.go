package http_handler_agents

import (
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_dto "github.com/sqlmerr/astragalaxy/internal/transport/http/dto"
	http_request "github.com/sqlmerr/astragalaxy/internal/transport/http/request"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

type CraftRequest struct {
	RecipeID          string    `json:"recipe_id" validate:"required"`
	TargetInventoryID uuid.UUID `json:"target_inventory_id" validate:"required"`
	Amount            int       `json:"amount" validate:"required,min=1"`
}

type CraftResponse struct {
	Cooldown http_dto.CooldownDTO `json:"cooldown"`
}

func (h *AgentsHTTPHandler) Craft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	var req CraftRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate HTTP request")
		return
	}

	agentID := core_auth.GetAgentIDFromContext(ctx)
	cooldown, err := h.craftingService.Craft(ctx, agentID, req.RecipeID, req.TargetInventoryID, req.Amount)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to process craft")
		return
	}

	response := CraftResponse{Cooldown: http_dto.ColdownFromModel(cooldown)}
	responseHandler.JSONResponse(http.StatusOK, response)
}
