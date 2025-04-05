package schema

type Planet struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SystemID string `json:"system_id"`
	Threat   string `json:"threat"`
}

type CreatePlanet struct {
	Name     string `json:"name"`
	SystemID string `json:"system_id"`
	Threat   string `json:"threat"`
}

type UpdatePlanet struct {
	Name     string `json:"name"`
	SystemID string `json:"system_id"`
	Threat   string `json:"threat"`
}
