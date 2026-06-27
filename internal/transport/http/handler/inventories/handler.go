package http_handler_inventories

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/service"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type InventoriesHTTPHandler struct {
	service service.Service
}

func NewInventoriesHTTPHandler(service service.Service) *InventoriesHTTPHandler {
	return &InventoriesHTTPHandler{service}
}

func (h *InventoriesHTTPHandler) Routes(agentAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/inventories/my",
			Handler:    h.GetMyInventory,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/inventories/my/ships/{id}",
			Handler:    h.GetMyShipInventory,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/inventories/transfer-resources",
			Handler:    h.TransferResources,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/inventories/transfer-items.yaml",
			Handler:    h.TransferItems,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}
