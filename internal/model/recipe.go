package model

import "time"

type RecipeResource struct {
	ResourceID string
	Amount     int
}

type Recipe struct {
	ID               string
	RequiredFacility FacilityType
	MinTier          int
	Duration         int
	Inputs           []RecipeResource
	Outputs          []RecipeResource
}

func (r *Recipe) GetDuration() time.Duration {
	return time.Duration(r.Duration) * time.Second
}
