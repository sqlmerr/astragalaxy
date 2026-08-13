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

type ItemDTO struct {
	ID               string `json:"id"`
	ProvidesFacility string `json:"provides_facility,omitempty"`
}

type ResourceDTO struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}

type FacilityDTO struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"`
	Tier           int     `json:"tier"`
	TimeMultiplier float64 `json:"time_multiplier"`
	CostMultiplier float64 `json:"cost_multiplier"`
}
