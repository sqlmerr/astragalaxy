package http_middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/auth"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	http_response "github.com/sqlmerr/astragalaxy/internal/transport/http/response"
)

const AGENT_ID_HEADER = "X-Agent-ID"

func UserAuth(jwtProcessor auth.JWTProcessor) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := http_response.NewHTTPResponseHandler(log, w)

			authorizationHeader := r.Header.Get("Authorization")
			headerParts := strings.SplitN(authorizationHeader, " ", 2)
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				responseHandler.ErrorResponse(
					fmt.Errorf("invalid jwt token: %w", errs.CodeInvalidJWTToken),
					"Invalid jwt token. Must be in format 'Bearer <jwt token>'",
				)
				return
			}

			jwtToken := headerParts[1]
			userID, err := jwtProcessor.ValidateToken(jwtToken)
			if err != nil {
				responseHandler.ErrorResponse(fmt.Errorf("%w: %w", err, errs.CodeInvalidJWTToken), "Invalid jwt token")
				return
			}
			ctx = context.WithValue(ctx, auth.UserIDContextKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type AgentGetter interface {
	GetAgentByToken(ctx context.Context, tokenHash string) (model.Agent, error)
	GetAgent(ctx context.Context, id uuid.UUID) (model.Agent, error)
}

func AgentAuth(jwtProcessor auth.JWTProcessor, agentGetter AgentGetter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := http_response.NewHTTPResponseHandler(log, w)

			authorizationHeader := r.Header.Get("Authorization")
			headerParts := strings.SplitN(authorizationHeader, " ", 2)
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				responseHandler.ErrorResponse(
					fmt.Errorf("invalid token: %w", errs.CodeInvalidAgentToken),
					"Invalid token. Must be in format 'Bearer <token>'",
				)
				return
			}

			var agent model.Agent
			if strings.HasPrefix(headerParts[1], "ag_agent_") {
				hashedToken := auth.HashRawAgentToken(headerParts[1])
				var err error
				agent, err = agentGetter.GetAgentByToken(ctx, hashedToken)
				if err != nil {
					if errors.Is(err, errs.ErrNotFound) {
						err = errs.CodeInvalidAgentToken
					}

					responseHandler.ErrorResponse(err, "Invalid agent token")
					return
				}
			} else {
				jwtToken := headerParts[1]
				userID, err := jwtProcessor.ValidateToken(jwtToken)
				if err != nil {
					responseHandler.ErrorResponse(fmt.Errorf("%w: %w", err, errs.CodeInvalidJWTToken), "Invalid jwt token")
					return
				}

				agentIDHeader := r.Header.Get(AGENT_ID_HEADER)
				agentID, err := uuid.Parse(agentIDHeader)
				if err != nil {
					responseHandler.ErrorResponse(
						fmt.Errorf("parse agent id: %w: %w", errs.CodeInvalidUUID, err),
						fmt.Sprintf("Failed to parse %s agent id header", AGENT_ID_HEADER),
					)
					return
				}

				agent, err = agentGetter.GetAgent(ctx, agentID)
				if err != nil {
					if errors.Is(err, errs.ErrNotFound) {
						err = errs.CodeAgentNotFound.WithBase(errs.ErrUnauthorized)
					}

					responseHandler.ErrorResponse(err, "Agent not found")
					return
				}
				if agent.UserID != userID {
					responseHandler.ErrorResponse(
						fmt.Errorf("agent with id=%s does not belong to user with id=%s: %w", agent.ID, userID, errs.CodeAccessDenied.WithBase(errs.ErrUnauthorized)),
						"Failed to access agent control",
					)
					return
				}
			}

			ctx = context.WithValue(ctx, auth.AgentIDContextKey, agent.ID)
			ctx = context.WithValue(ctx, auth.AgentContextKey, agent)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
