package model

import "time"

type LocationType string

const (
	LocationWaypoint LocationType = "WAYPOINT"
	LocationPlanet   LocationType = "PLANET"
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
