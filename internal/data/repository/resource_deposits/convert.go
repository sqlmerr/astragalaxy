package resource_deposits_repository

import (
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	"github.com/sqlmerr/astragalaxy/internal/model"
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
