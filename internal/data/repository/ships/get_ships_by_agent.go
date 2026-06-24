package ships_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

func (r *ShipRepositoryImpl) GetShipsByAgent(ctx context.Context, agentID uuid.UUID) ([]model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	query := `
	SELECT id, agent_id, type, active, system_x, system_y, status, created_at, name
	FROM ships
	WHERE agent_id = $1
	ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query, agentID)
	if err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	defer rows.Close()

	var ships []model.Ship
	for rows.Next() {
		var s model.Ship
		err := rows.Scan(&s.ID, &s.AgentID, &s.Type, &s.Active, &s.SystemX, &s.SystemY, &s.Status, &s.CreatedAt, &s.Name)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		ships = append(ships, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return ships, nil
}
