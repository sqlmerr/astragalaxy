package http_handler_agents

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type AgentResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func agentDTOFromModel(m model.Agent) AgentResponseDTO {
	return AgentResponseDTO{
		ID:        m.ID,
		UserID:    m.UserID,
		Username:  m.Username,
		CreatedAt: m.CreatedAt,
	}
}

func agentDTOsFromModels(m []model.Agent) []AgentResponseDTO {
	return lo.Map(m, func(i model.Agent, _ int) AgentResponseDTO {
		return agentDTOFromModel(i)
	})
}

type CooldownDTO struct {
	SetAt    time.Time     `json:"set_at"`
	Duration time.Duration `json:"duration"`
	Action   string        `json:"action"`
}
