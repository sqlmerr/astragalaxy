package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *ShipRepositoryImpl) CreateShip(ctx context.Context, data CreateShip) (model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO ships (agent_id, type, active, system_x, system_y, status, name) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, agent_id, type, active, system_x, system_y, status, name;
	`

	row := r.db.QueryRow(ctx, query, data.AgentID, data.Type, data.Active, data.SystemX, data.SystemY, data.Status, data.Name)
	var s model.Ship
	err := row.Scan(&s.ID, &s.AgentID, &s.Type, &s.Active, &s.SystemX, &s.SystemY, &s.Status, &s.Name)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.Ship{}, core_errors.NewWithCode(
				core_errors.CodeAgentNotFound,
				fmt.Errorf("agent with id='%s': %w", data.AgentID, core_errors.ErrNotFound),
			)
		}

		return model.Ship{}, fmt.Errorf("scan: %w", err)
	}

	return s, nil
}
