package model

import "github.com/lib/pq"

type System struct {
	ID          string             `gorm:"not null;primaryKey"`
	Name        string             `gorm:"not null"`
	Connections []SystemConnection `gorm:"foreignKey:SystemFromID;constraint:OnDelete:CASCADE"`
	Locations pq.StringArray `gorm:"type:text[]"`
}
