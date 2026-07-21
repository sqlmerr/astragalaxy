package ships_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
)

func (r *ShipRepositoryImpl) CreateShip(ctx context.Context, data CreateShip) (model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	s, err := r.q.CreateShip(ctx, database.CreateShipParams{
		AgentID:     data.AgentID,
		Type:        database.ShipType(data.Type),
		Active:      data.Active,
		SystemX:     int32(data.SystemX),
		SystemY:     int32(data.SystemY),
		Status:      database.ShipStatus(data.Status),
		Name:        data.Name,
		InventoryID: data.InventoryID,
		Location:    database.ShipLocation(data.Location),
		LocationID:  int32(data.LocationID),
	})
	err = postgres_pool.TranslateError(err)

	if err != nil {
		if errors.Is(err, postgres_pool.ErrViolatesForeignKey) {
			return model.Ship{}, core_errors.NewWithCode(
				core_errors.CodeAgentNotFound,
				fmt.Errorf("agent with id='%s': %w", data.AgentID, core_errors.ErrNotFound),
			)
		}

		return model.Ship{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(s), nil
}
