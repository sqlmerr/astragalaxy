package item_actions

import (
	"context"

	"github.com/sqlmerr/astragalaxy/internal/data"
)

type ActionContext struct {
	Context context.Context
	Store   data.Store
}
