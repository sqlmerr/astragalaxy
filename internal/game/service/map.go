package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
)

func (s *Service) ShipRadar(ctx context.Context, agentID uuid.UUID) ([]worldgen.System, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("get active ship: %w", err)
	}

	systems, err := s.worldGen.GetSystemsInBox(ship.SystemX-10, ship.SystemY-10, ship.SystemX+10, ship.SystemY+10)
	if err != nil {
		return nil, fmt.Errorf("use radar: %w", err)
	}

	return systems, nil
}

func (s *Service) GetCurrentAgentSystem(ctx context.Context, agentID uuid.UUID) (worldgen.System, error) {
	ship, err := s.store.Ships().GetActiveShipByAgent(ctx, agentID)
	if err != nil {
		return worldgen.System{}, fmt.Errorf("get active ship: %w", err)
	}

	system, exists := s.worldGen.GenerateSystemByCoords(ship.SystemX, ship.SystemY)
	if !exists {
		return worldgen.System{}, core_errors.NewWithCode(
			core_errors.CodeAnomaly,
			fmt.Errorf(
				"something happened to system with x=%d y=%d: %w",
				ship.SystemX,
				ship.SystemY,
				core_errors.ErrUnprocessableEntity,
			),
		)
	}

	return *system, nil
}
