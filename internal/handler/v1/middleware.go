package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Middleware = func(ctx huma.Context, next func(huma.Context))

func (h *Handler) SudoMiddleware(api huma.API) Middleware {
	return func(ctx huma.Context, next func(huma.Context)) {
		header := ctx.Header("secret-token")
		if header != h.state.Config.SecretToken {
			util.WriteError(api, ctx, util.ErrInvalidToken)
			return
		}

		next(ctx)
	}
}

func (h *Handler) UserGetter(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		username := ctx.Context().Value("subject").(string)

		user, err := h.s.FindOneUserByUsername(username)
		if err != nil || user == nil {
			util.WriteError(api, ctx, util.ErrInvalidToken)
			return
		}

		ctx = huma.WithValue(ctx, "user", user)

		next(ctx)
	}
}

func (h *Handler) AstralGetter(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user := ctx.Context().Value("user").(*schema.User)

		astralID := ctx.Header("X-Astral-ID")
		baseErr := util.ErrInvalidAstralIDHeader
		if astralID == "" {
			util.WriteError(api, ctx, baseErr)
			return
		}

		ID, err := uuid.Parse(astralID)
		if err != nil {
			util.WriteError(api, ctx, baseErr)
			return
		}

		astral, err := h.s.FindOneAstral(ID)
		if err != nil || astral == nil {
			util.WriteError(api, ctx, baseErr)
			return
		}

		if astral.UserID != user.ID {
			util.WriteError(api, ctx, baseErr)
			return
		}

		ctx = huma.WithValue(ctx, "astral", astral)
		next(ctx)
	}
}

func (h *Handler) JWTMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
		if token == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Missing token")
			return
		}
		parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.state.Config.JwtSecret), nil
		})
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Malformed token")
			return
		}
		sub, err := parsed.Claims.GetSubject()
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Malformed token")
			return
		}
		ctx = huma.WithValue(ctx, "subject", sub)
		next(ctx)
	}
}
