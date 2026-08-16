package ships_repository

import (
	"context"
	"errors"
	"fmt"

	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func (r *ShipRepositoryImpl) SaveShip(ctx context.Context, ship model.Ship) (model.Ship, error) {
	ctx, cancel := context.WithTimeout(ctx, r.db.OpTimeout())
	defer cancel()

	s, err := r.q.SaveShip(ctx, database.SaveShipParams{
		ID:         ship.ID,
		Type:       database.ShipType(ship.Type),
		Active:     ship.Active,
		SystemX:    int32(ship.Coords.SystemX),
		SystemY:    int32(ship.Coords.SystemY),
		Status:     database.ShipStatus(ship.Status),
		Name:       ship.Name,
		Location:   database.ShipLocation(ship.Coords.Location),
		LocationID: int32(ship.Coords.LocationID),
	})
	err = postgres_pool.TranslateError(err)
	if err != nil {
		if errors.Is(err, postgres_pool.ErrNoRows) {
			return model.Ship{}, errs.NewWithCode(
				errs.CodeShipNotFound,
				fmt.Errorf("ship with id='%s': %w", ship.ID, errs.ErrNotFound),
			)
		}

		return model.Ship{}, fmt.Errorf("scan: %w", err)
	}

	return convertModel(s), nil
}
