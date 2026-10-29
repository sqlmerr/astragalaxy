package http_handler_wallets

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/wallets"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type WalletsHTTPHandler struct {
	service wallets.Service
}

func New(service wallets.Service) *WalletsHTTPHandler {
	return &WalletsHTTPHandler{service}
}

func (h *WalletsHTTPHandler) Routes(agentAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/wallets/my",
			Handler:    h.GetMyWallets,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/wallets/my/{addr}/pay",
			Handler:    h.Pay,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/wallets",
			Handler:    h.RegisterWallet,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}
