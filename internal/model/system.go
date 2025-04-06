package model

type System struct {
	ID          string             `gorm:"not null;primaryKey"`
	Name        string             `gorm:"not null"`
	Connections []SystemConnection `gorm:"foreignKey:SystemFromID"`
}
