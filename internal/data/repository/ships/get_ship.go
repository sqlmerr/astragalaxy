package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *ShipRepositoryImpl) GetShip(ctx context.Context, shipID uuid.UUID) (model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	query := `
	SELECT id, agent_id, type, active, system_x, system_y, status, created_at, name
	FROM ships
	WHERE id = $1;
	`

	row := r.db.QueryRow(ctx, query, shipID)
	var s model.Ship
	err := row.Scan(
		&s.ID,
		&s.AgentID,
		&s.Type,
		&s.Active,
		&s.SystemX,
		&s.SystemY,
		&s.Status,
		&s.CreatedAt,
		&s.Name,
	)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Ship{}, core_errors.NewWithCode(
				core_errors.CodeShipNotFound,
				fmt.Errorf("ship with id='%s': %w", shipID, core_errors.ErrNotFound),
			)
		}

		return model.Ship{}, fmt.Errorf("scan: %w", err)
	}

	return s, nil
}
