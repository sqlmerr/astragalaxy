package schema

import "github.com/google/uuid"

type Inventory struct {
	ID       uuid.UUID `json:"id"`
	HolderID uuid.UUID `json:"holder_id"`
	Holder   string    `json:"holder_type"`
	Items    []Item    `json:"items"`
}
