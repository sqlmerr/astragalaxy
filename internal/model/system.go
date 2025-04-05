package model

type System struct {
	ID   string `gorm:"not null"`
	Name string `gorm:"not null"`
}
