package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *ShipRepositoryImpl) SaveShip(ctx context.Context, ship model.Ship) (model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	query := `
	UPDATE ships 
	SET 
	    type = $2,
	    active = $3,
	    system_x = $4,
	    system_y = $5,
	    status = $6,
	    name = $7
	WHERE id = $1
	RETURNING id, agent_id, type, active, system_x, system_y, status, created_at, name
	`

	row := r.db.QueryRow(
		ctx,
		query,
		ship.ID,
		ship.Type,
		ship.Active,
		ship.SystemX,
		ship.SystemY,
		ship.Status,
		ship.Name,
	)

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
				fmt.Errorf("ship with id='%s': %w", ship.ID, core_errors.ErrNotFound),
			)
		}

		return model.Ship{}, fmt.Errorf("scan: %w", err)
	}

	return s, nil
}
