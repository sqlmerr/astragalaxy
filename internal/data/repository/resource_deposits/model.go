package resource_deposits_repository

import (
	"time"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
)

type CreateResourceDeposit struct {
	SystemX      int
	SystemY      int
	LocationType model.LocationType
	LocationID   int
	ResourceType model.ResourceType
	Remaining    int
	LastMinedAt  time.Time
}

type GetResourceDeposit struct {
	SystemX      int
	SystemY      int
	LocationType model.LocationType
	LocationID   int
	ResourceType model.ResourceType
}
