package http_handler_agents

import (
	"net/http"

	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func (h *AgentsHTTPHandler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	// TODO
	responseHandler.JSONResponse(200, "placeholder")
}
