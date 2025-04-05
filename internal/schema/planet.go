package schema

type PlanetSchema struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SystemID string `json:"system_id"`
	Threat   string `json:"threat"`
}

type CreatePlanetSchema struct {
	Name     string `json:"name"`
	SystemID string `json:"system_id"`
	Threat   string `json:"threat"`
}

type UpdatePlanetSchema struct {
	Name     string `json:"name"`
	SystemID string `json:"system_id"`
	Threat   string `json:"threat"`
}
