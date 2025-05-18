package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type ResourcesJSON map[string]int

func (a ResourcesJSON) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ResourcesJSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type Bundle struct {
	ID        uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Resources ResourcesJSON `gorm:"type:jsonb"`

	InventoryID uuid.UUID `gorm:"not null"`
	Inventory   Inventory
}
