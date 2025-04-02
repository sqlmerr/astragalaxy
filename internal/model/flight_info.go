package model

import "github.com/google/uuid"

type FlightInfo struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Flying        *bool     `gorm:"not null;default:false"`
	Destination   string    // planet or system
	DestinationID uuid.UUID // planet or system id
	FlownOutAt    int64
	FlyingTime    int64
}
