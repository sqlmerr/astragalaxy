package model

import "time"

type RecipeInput struct {
	ResourceID string
	Amount     int
}

type RecipeOutputType string

const (
	RecipeOutputItem     RecipeOutputType = "item"
	RecipeOutputResource RecipeOutputType = "resource"
)

type RecipeOutput struct {
	Type   RecipeOutputType
	ID     string
	Amount int
}

type Recipe struct {
	ID               string
	RequiredFacility FacilityType
	MinTier          int
	Duration         int
	Inputs           []RecipeInput
	Outputs          []RecipeOutput
}

func (r *Recipe) GetDuration() time.Duration {
	return time.Duration(r.Duration) * time.Second
}
