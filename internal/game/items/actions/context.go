package item_actions

import (
	"context"

	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
)

type ActionContext struct {
	Context  context.Context
	Store    data.Store
	GameData registry.GameData
}
