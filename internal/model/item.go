package model

import "github.com/google/uuid"

type Item struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InventoryID uuid.UUID
	Inventory   Inventory
	Code        string
	Durability  int `gorm:"default:100"`
}
