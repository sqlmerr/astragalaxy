package schema

import (
	"astragalaxy/internal/model"

	"github.com/google/uuid"
)

type ExplorationInfo struct {
	ID           uuid.UUID             `json:"id"`
	Status       bool                  `json:"status"`
	Type         model.ExplorationType `json:"type"`
	StartedAt    int64                 `json:"started_at"`
	RequiredTime int64                 `json:"required_time"`
}

func ExplorationInfoSchemaFromExplorationInfo(info *model.ExplorationInfo) ExplorationInfo {
	return ExplorationInfo{
		ID:           info.ID,
		Status:       info.Exploring,
		Type:         info.Type,
		StartedAt:    info.StartedAt,
		RequiredTime: info.RequiredTime,
	}
}
