package http_handler_data

type RecipeResourceDTO struct {
	ResourceID string `json:"resource_id"`
	Amount     int    `json:"amount"`
}

type RecipeDTO struct {
	ID               string              `json:"id"`
	RequiredFacility string              `json:"required_facility"`
	MinTier          int                 `json:"min_tier"`
	Duration         int                 `json:"duration"`
	Inputs           []RecipeResourceDTO `json:"inputs"`
	Outputs          []RecipeResourceDTO `json:"outputs"`
}
