package wallets

import (
	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type PayInput struct {
	AgentID     uuid.UUID
	FromAddress string
	ToAddress   string
	Units       int
	Quarks      int
}

type PayOutput struct {
	Fee      int // in units
	Cooldown model.Cooldown
}
