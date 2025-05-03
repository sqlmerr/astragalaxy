package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) registerWalletGroup(api huma.API) {
	tags := []string{"wallets"}
	security := []map[string][]string{{"bearerAuth": {}}}
	params := []*huma.Param{{Name: "X-Astral-ID", In: "header", Description: "astral id", Required: true, Schema: &huma.Schema{Type: "string"}}}
	api.UseMiddleware(h.JWTMiddleware(api), h.UserGetter(api), h.AstralGetter(api))

	huma.Register(api, huma.Operation{
		Method:     http.MethodGet,
		Path:       "/my",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getMyWallets)
	huma.Register(api, huma.Operation{
		Method:     http.MethodPost,
		Path:       "/my",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.createMyWallet)
	huma.Register(api, huma.Operation{
		Method:     http.MethodPost,
		Path:       "/my/{id}/send",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.sendCurrency)
	huma.Register(api, huma.Operation{
		Method:     http.MethodPatch,
		Path:       "/my/{id}/lock",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.lockMyWallet)
	huma.Register(api, huma.Operation{
		Method:     http.MethodPatch,
		Path:       "/my/{id}/unlock",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.unlockMyWallet)
	huma.Register(api, huma.Operation{
		Method:     http.MethodGet,
		Path:       "/{id}",
		Tags:       tags,
		Security:   security,
		Parameters: params,
	}, h.getWalletByID)
}

func (h *Handler) getWalletByID(_ context.Context, input *struct {
	ID string `path:"id"`
}) (*schema.BaseResponse[schema.Wallet], error) {
	id, err := uuid.Parse(input.ID)
	if err != nil || id == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	wallet, err := h.s.GetWallet(id)
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.Wallet]{Body: *wallet}, nil
}

func (h *Handler) getMyWallets(ctx context.Context, _ *struct{}) (*schema.BaseDataResponse[[]schema.Wallet], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	wallets, err := h.s.GetAstralWallets(astral.ID)
	if err != nil {
		return nil, err
	}
	response := schema.BaseDataResponse[[]schema.Wallet]{
		Body: schema.DataGenericResponse[[]schema.Wallet]{
			Data: wallets,
		},
	}
	return &response, nil
}

func (h *Handler) createMyWallet(ctx context.Context, input *schema.BaseRequest[schema.CreateWallet]) (*schema.BaseResponse[schema.Wallet], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	wallets, err := h.s.GetAstralWallets(astral.ID)
	if err != nil {
		return nil, err
	}
	if len(wallets) == 3 {
		return nil, util.New("this astral has reached maximum of wallets (3)", http.StatusBadRequest)
	}

	wallet, err := h.s.CreateWallet(input.Body, astral.ID)
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.Wallet]{Body: *wallet}, nil
}

func (h *Handler) sendCurrency(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body schema.WalletTransaction
}) (*schema.BaseResponse[schema.OkResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	id, err := uuid.Parse(input.ID)
	if err != nil || id == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	wallet, err := h.s.GetWallet(id)
	if err != nil {
		return nil, err
	}

	if wallet.AstralID != astral.ID {
		return nil, util.ErrNotFound
	}

	err = h.s.ProceedTransaction(wallet.ID, &input.Body)
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}

func (h *Handler) lockMyWallet(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*schema.BaseResponse[schema.OkResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	id, err := uuid.Parse(input.ID)
	if err != nil || id == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	wallet, err := h.s.GetWallet(id)
	if err != nil {
		return nil, err
	}

	if wallet.AstralID != astral.ID {
		return nil, util.ErrNotFound
	}

	err = h.s.LockWallet(wallet.ID)
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}

func (h *Handler) unlockMyWallet(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*schema.BaseResponse[schema.OkResponse], error) {
	astral := ctx.Value("astral").(*schema.Astral)
	id, err := uuid.Parse(input.ID)
	if err != nil || id == uuid.Nil {
		return nil, util.ErrIDMustBeUUID
	}

	wallet, err := h.s.GetWallet(id)
	if err != nil {
		return nil, err
	}

	if wallet.AstralID != astral.ID {
		return nil, util.ErrNotFound
	}

	err = h.s.UnlockWallet(wallet.ID)
	if err != nil {
		return nil, err
	}

	return &schema.BaseResponse[schema.OkResponse]{Body: schema.OkResponse{Ok: true, CustomStatusCode: 1}}, nil
}
