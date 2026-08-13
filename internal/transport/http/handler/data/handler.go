package http_handler_data

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type DataHTTPHandler struct {
	gameData registry.GameData
}

func New(gameData registry.GameData) *DataHTTPHandler {
	return &DataHTTPHandler{gameData}
}

func (h *DataHTTPHandler) Routes() []http_server.Route {
	return []http_server.Route{
		{
			Path:    "/data/recipes",
			Method:  http.MethodGet,
			Handler: h.GetRecipes,
		},
		{
			Path:    "/data/items",
			Method:  http.MethodGet,
			Handler: h.GetItems,
		},
		{
			Path:    "/data/resources",
			Method:  http.MethodGet,
			Handler: h.GetResources,
		},
		{
			Path:    "/data/facilities",
			Method:  http.MethodGet,
			Handler: h.GetFacilities,
		},
	}
}
