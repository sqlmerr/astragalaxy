package model

import "time"

type LocationType string

const (
	LocationWaypoint LocationType = "waypoint"
	LocationPlanet   LocationType = "planet"

	ResourceDepositRegenerationRate int = 1 // amount of resources that will be regenerated every minute
)

type ResourceDeposit struct {
	SystemX      int
	SystemY      int
	LocationType LocationType
	LocationID   int
	ResourceType ResourceType
	Remaining    int
	LastMinedAt  time.Time
}

func (d *ResourceDeposit) Regenerate(maxResourceAmount int) {
	if d.Remaining >= maxResourceAmount {
		return
	}
	now := time.Now()
	elapsed := now.Sub(d.LastMinedAt)

	d.Remaining += int(elapsed.Minutes()) * ResourceDepositRegenerationRate
	if d.Remaining > maxResourceAmount {
		d.Remaining = maxResourceAmount
	}
}
