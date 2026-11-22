package http_handler_mail

import (
	"net/http"

	"github.com/sqlmerr/astragalaxy/internal/game/mail"
	http_middleware "github.com/sqlmerr/astragalaxy/internal/transport/http/middleware"
	http_server "github.com/sqlmerr/astragalaxy/internal/transport/http/server"
)

type MailHTTPHandler struct {
	service mail.Service
}

func New(service mail.Service) *MailHTTPHandler {
	return &MailHTTPHandler{service: service}
}

func (h *MailHTTPHandler) Routes(agentAuthMiddleware http_middleware.Middleware) []http_server.Route {
	return []http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/mail/my",
			Handler:    h.GetClaimablePackages,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
		{
			Method:     http.MethodPost,
			Path:       "/mail/{id}/claim",
			Handler:    h.ClaimPackage,
			Middleware: []http_middleware.Middleware{agentAuthMiddleware},
		},
	}
}
