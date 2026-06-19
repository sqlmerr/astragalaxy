package http_handler_users

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type UsersHTTPHandler struct {
	service game.Service
}

func NewUsersHTTPHandler(service game.Service) *UsersHTTPHandler {
	return &UsersHTTPHandler{service}
}

func (h *UsersHTTPHandler) Routes(userAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.RegisterUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.LoginUser,
		},
		{
			Method:     http.MethodGet,
			Path:       "/auth/me",
			Handler:    h.GetMe,
			Middleware: []http_middleware.Middleware{userAuthMiddleware},
		},
	}
}
