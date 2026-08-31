package model

import "github.com/google/uuid"

type ResourceType string

const (
	ResourceWarpCellT1 ResourceType = "warp_cell_t1"
	ResourceWarpCellT2 ResourceType = "warp_cell_t2"
)

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
