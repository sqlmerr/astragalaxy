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

type FacilityCraftingModifiers struct {
	TimeMultiplier float64
	CostMultiplier float64
}

type Facility struct {
	ID   string
	Type FacilityType
	Tier int

	// FacilityCraftingModifiers or ...
	Modifiers any
}

func (f *Facility) CraftingMods() (*FacilityCraftingModifiers, bool) {
	m, ok := f.Modifiers.(*FacilityCraftingModifiers)
	return m, ok
}
