package item_actions

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type ActionResponse interface {
	Response() json.RawMessage
	Cooldown() time.Duration
}

type okActionResponse struct {
	cooldownDuration time.Duration
}

func (r okActionResponse) Response() json.RawMessage {
	return json.RawMessage(`{"ok": true}`)
}

func (r okActionResponse) Cooldown() time.Duration {
	return r.cooldownDuration
}

func newOkResponse(cooldownDuration time.Duration) okActionResponse {
	return okActionResponse{cooldownDuration: cooldownDuration}
}

type ItemAction func(ctx ActionContext, agentID uuid.UUID, item model.Item) (ActionResponse, error)

var Actions = map[model.ItemType]ItemAction{
	model.ItemPortableSmelter: UsePortableSmelter,
}
