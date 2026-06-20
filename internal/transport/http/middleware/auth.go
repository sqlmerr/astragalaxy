package http_middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	core_auth "github.com/sqlmerr/astragalaxy/internal/auth"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

func UserAuth(jwtProcessor core_auth.JWTProcessor) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := http_response.NewHTTPResponseHandler(log, w)

			authorizationHeader := r.Header.Get("Authorization")
			headerParts := strings.SplitN(authorizationHeader, " ", 2)
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				responseHandler.ErrorResponse(
					core_errors.NewWithCode(core_errors.CodeInvalidJWTToken, fmt.Errorf("invalid jwt token: %w", core_errors.ErrUnauthorized)),
					"Invalid jwt token. Must be in format 'Bearer <jwt token>'",
				)
				return
			}

			jwtToken := headerParts[1]
			userID, err := jwtProcessor.ValidateToken(jwtToken)
			if err != nil {
				responseHandler.ErrorResponse(core_errors.NewWithCode(core_errors.CodeInvalidJWTToken, err), "Invalid jwt token")
				return
			}
			ctx = context.WithValue(ctx, core_auth.UserIDContextKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type AgentGetter interface {
	GetAgentByToken(ctx context.Context, tokenHash string) (model.Agent, error)
}

func AgentAuth(agentGetter AgentGetter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := http_response.NewHTTPResponseHandler(log, w)

			authorizationHeader := r.Header.Get("Authorization")
			headerParts := strings.SplitN(authorizationHeader, " ", 2)
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				responseHandler.ErrorResponse(
					core_errors.NewWithCode(core_errors.CodeInvalidAgentToken, fmt.Errorf("invalid agent token: %w", core_errors.ErrUnauthorized)),
					"Invalid agent token. Must be in format 'Bearer ag_agent_<token>'",
				)
				return
			}
			hashedToken := core_auth.HashRawAgentToken(headerParts[1])
			agent, err := agentGetter.GetAgentByToken(ctx, hashedToken)
			if err != nil {
				if errors.Is(err, core_errors.ErrNotFound) {
					err = core_errors.NewWithCode(core_errors.CodeInvalidAgentToken, core_errors.ErrUnauthorized)
				}

				responseHandler.ErrorResponse(err, "Invalid agent token")
				return
			}

			ctx = context.WithValue(ctx, core_auth.AgentIDContextKey, agent.ID)
			ctx = context.WithValue(ctx, core_auth.AgentContextKey, agent)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
