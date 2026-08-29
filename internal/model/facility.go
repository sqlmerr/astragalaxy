package model

import (
	"fmt"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
)

type FacilityType string

const (
	FacilityPrinter FacilityType = "printer"
	FacilitySmelter FacilityType = "smelter"
)

func NewFacilityType(facilityType string) (FacilityType, error) {
	f := FacilityType(facilityType)
	switch f {
	case FacilityPrinter, FacilitySmelter:
		return f, nil
	default:
		return "", errs.NewWithCode(errs.CodeInvalidFacilityType, fmt.Errorf("invalid facility type %s: %w", facilityType, errs.ErrInvalidArgument))
	}
}

type Facility struct {
	ID             string
	Type           FacilityType
	Tier           int
	TimeMultiplier float64
	CostMultiplier float64
}
