package model

import "github.com/google/uuid"

type ResourceType string

// const (
// 	// basic
// 	ResourceIron    ResourceType = "iron"
// 	ResourceCrystal ResourceType = "crystal"
// 	ResourceCarbon  ResourceType = "carbon"
// 	ResourceIce     ResourceType = "ice"

// 	// advanced
// 	ResourceCopper   ResourceType = "copper"
// 	ResourceTitanium ResourceType = "titanium"
// 	ResourceSilicon  ResourceType = "silicon"
// 	ResourceHelium   ResourceType = "helium"

// 	// exotic
// 	ResourceUranium     ResourceType = "uranium"
// 	ResourceIridium     ResourceType = "iridium"
// 	ResourceDarkMatter  ResourceType = "dark_matter"
// 	ResourceBioDisputes ResourceType = "bio_disputes"

// 	// composite
// 	ResourceSteel ResourceType = "steel"
// )

type Resource struct {
	InventoryID  uuid.UUID
	ResourceType ResourceType
	Amount       int
}

type ResourceDataWorldGenParams struct {
	Min int
	Max int
}

type ResourceData struct {
	ID       string
	Tags     []string
	WorldGen map[string]ResourceDataWorldGenParams
}
