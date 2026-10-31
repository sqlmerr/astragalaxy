package item_actions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type UseContractInput struct {
	Address string `json:"address"`
}

type ContactMetadata struct {
	Units  int `json:"units,omitempty"`
	Quarks int `json:"quarks,omitempty"`
}

type useContactResponse struct {
	Units  int
	Quarks int
}

func (r useContactResponse) Response() json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{"ok": true, "units": %d, "quarks": %d}`, r.Units, r.Quarks))
}

func (r useContactResponse) Cooldown() time.Duration {
	return 5 * time.Second
}

func UseContract(ctx ActionContext, agentID uuid.UUID, item model.Item, input json.RawMessage) (ActionResponse, error) {
	var i UseContractInput
	if err := json.Unmarshal(input, &i); err != nil {
		return nil, errs.NewWithCode(errs.CodeDecodeError, fmt.Errorf("decode action input json: %w: %w", err, errs.ErrInvalidArgument))
	}

	wallet, err := ctx.Store.Wallets().GetWallet(ctx.Context, i.Address)
	if err != nil {
		return nil, fmt.Errorf("get wallet by address = %q: %w", i.Address, err)
	}

	if wallet.OwnerID != agentID {
		return nil, errs.NewWithCode(errs.CodeAccessDenied, fmt.Errorf("cannot access this wallet (%s): %w", i.Address, errs.ErrAccessDenied))
	}

	var m ContactMetadata
	if err := json.Unmarshal(item.Metadata, &m); err != nil {
		return nil, errs.NewWithCode(errs.CodeDecodeError, fmt.Errorf("decode item metadata json: %w", err))
	}

	wallet.Units += m.Units
	wallet.Quarks += m.Quarks
	_, err = ctx.Store.Wallets().SaveWallet(ctx.Context, wallet)
	if err != nil {
		return nil, fmt.Errorf("save wallet: %w", err)
	}

	err = ctx.Store.Inventories().DeleteItem(ctx.Context, item.ID)
	if err != nil {
		return nil, fmt.Errorf("delete item: %w", err)
	}

	return useContactResponse{Units: m.Units, Quarks: m.Quarks}, nil
}
