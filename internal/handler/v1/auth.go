package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) registerAuthGroup(api huma.API) {
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/register", DefaultStatus: 201, Tags: []string{"auth"}}, h.registerUser)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/login", Tags: []string{"auth"}}, h.login)
	huma.Register(api, huma.Operation{Method: http.MethodPost, Path: "/login/token", Tags: []string{"auth"}}, h.loginByToken)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/me",
		Description: "Get authorized user",
		Middlewares: huma.Middlewares{h.JWTMiddleware(api), h.UserGetter(api)},
		Security:    []map[string][]string{{"bearerAuth": {}}},
		Tags:        []string{"auth"},
	}, h.getMe)
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/token/sudo",
		Description: "Get user token sudo",
		Middlewares: huma.Middlewares{h.SudoMiddleware(api)},
		Security:    []map[string][]string{{"sudoAuth": {}}},
		Tags:        []string{"auth"},
	}, h.getUserTokenSudo)
}

func (h *Handler) registerUser(_ context.Context, input *schema.BaseRequest[schema.CreateUser]) (*schema.BaseResponse[schema.User], error) {
	system, err := h.s.FindOneSystemByName("initial")
	if err != nil || system == nil {
		return nil, err
	}

	user, err := h.s.Register(input.Body)
	if err != nil || user == nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.User]{Body: *user, Status: http.StatusCreated}, nil
}

func (h *Handler) loginByToken(_ context.Context, input *schema.BaseRequest[schema.AuthPayloadToken]) (*schema.BaseResponse[schema.AuthBody], error) {
	jwtToken, err := h.s.LoginByToken(input.Body.Token)

	if err != nil || jwtToken == nil {
		return nil, util.ErrUnauthorized
	}

	return &schema.BaseResponse[schema.AuthBody]{Body: schema.AuthBody{AccessToken: *jwtToken, TokenType: "Bearer"}}, nil
}

func (h *Handler) login(_ context.Context, input *schema.BaseRequest[schema.AuthPayload]) (*schema.BaseResponse[schema.AuthBody], error) {
	jwtToken, err := h.s.Login(&input.Body)

	if err != nil || jwtToken == nil {
		return nil, util.ErrUnauthorized
	}

	return &schema.BaseResponse[schema.AuthBody]{Body: schema.AuthBody{AccessToken: *jwtToken, TokenType: "Bearer"}}, nil
}

func (h *Handler) getMe(ctx context.Context, _ *struct{}) (*schema.BaseResponse[schema.User], error) {
	user := ctx.Value("user").(*schema.User)
	return &schema.BaseResponse[schema.User]{Body: *user}, nil
}

func (h *Handler) getUserTokenSudo(_ context.Context, input *struct {
	UserID string `query:"id" doc:"User id in UUID type" example:"39518366-f039-4a79-961f-611f6a2fe723"`
}) (*schema.BaseResponse[schema.UserTokenResponse], error) {
	ID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("failed to parse user id", err)
	}

	user, err := h.s.FindOneUserRaw(ID)
	if err != nil || user == nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.UserTokenResponse]{Body: schema.UserTokenResponse{Token: user.Token}}, nil
}
