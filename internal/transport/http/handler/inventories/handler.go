package http_handler_inventories

import (
	"net/http"

	inventory_service "github.com/sqlmerr/astragalaxy/internal/game/inventory"
	items_service "github.com/sqlmerr/astragalaxy/internal/game/items"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type InventoriesHTTPHandler struct {
	inventoryService inventory_service.InventoryService
	itemsService     items_service.ItemsService
}

func NewInventoriesHTTPHandler(inventoryService inventory_service.InventoryService, itemsService items_service.ItemsService) *InventoriesHTTPHandler {
	return &InventoriesHTTPHandler{inventoryService, itemsService}
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
			Path:       "/inventories/transfer-items",
			Handler:    h.TransferItems,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/inventories/my/items/{id}/use",
			Handler:    h.UseItem,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}
