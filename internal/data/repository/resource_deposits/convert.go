package resource_deposits_repository

import (
	"github.com/sqlmerr/astragalaxy/internal/data/model"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
)

func resourceDepositModel(deposit database.ResourceDepositState) model.ResourceDeposit {
	return model.ResourceDeposit{
		SystemX: int(deposit.SystemX), SystemY: int(deposit.SystemY),
		LocationType: model.LocationType(deposit.LocType),
		LocationID:   int(deposit.LocID),
		ResourceType: model.ResourceType(deposit.ResourceType),
		Remaining:    int(deposit.Remaining),
		LastMinedAt:  deposit.LastMinedAt.Time,
	}
}
