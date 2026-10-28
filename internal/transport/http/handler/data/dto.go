package http_handler_data

type RecipeInputDTO struct {
	ResourceID string `json:"resource_id"`
	Amount     int    `json:"amount"`
}

type RecipeOutputDTO struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}

type RecipeDTO struct {
	ID               string            `json:"id"`
	RequiredFacility string            `json:"required_facility"`
	MinTier          int               `json:"min_tier"`
	Duration         int               `json:"duration"`
	Inputs           []RecipeInputDTO  `json:"inputs"`
	Outputs          []RecipeOutputDTO `json:"outputs"`
}

type ItemShipModuleMining struct {
	SpeedMultiplier float64 `json:"speed_multiplier"`
}

type ItemShipModule struct {
	ProvidesFacility *string               `json:"provides_facility"`
	Mining           *ItemShipModuleMining `json:"mining"`
}

type ItemDTO struct {
	ID         string          `json:"id"`
	ShipModule *ItemShipModule `json:"ship_module,omitempty"`
}

type ResourceDTO struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}

type FacilityCraftingModsDTO struct {
	TimeMultiplier float64 `json:"time_multiplier"`
	CostMultiplier float64 `json:"cost_multiplier"`
}

type FacilityDTO struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Tier      int    `json:"tier"`
	Modifiers any    `json:"modifiers"`
}

type ShipClassDTO struct {
	ID string `json:"id"`

	MaxModuleSlots      int     `json:"max_module_slots"`
	MaxWarpDistance     int     `json:"max_warp_distance"`
	MiningLimitModifier float64 `json:"mining_limit_modifier"`
}
