package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
)

type FacilityType string

const (
	FacilityPrinter FacilityType = "printer"
	FacilitySmelter FacilityType = "smelter"
)

type Facility struct {
	ID             string       `json:"id"`
	Type           FacilityType `json:"type"`
	Tier           int          `json:"tier"`
	TimeMultiplier float64      `json:"time_multiplier"`
	CostMultiplier float64      `json:"cost_multiplier"`
}

type FacilityRegistry struct {
	facilities []Facility
}

func NewFacilityRegistry() *FacilityRegistry {
	return &FacilityRegistry{facilities: nil}
}

func (r *FacilityRegistry) Load(cfg Config) error {
	file, err := os.ReadFile(cfg.FacilitiesPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	var facilities []Facility
	err = json.Unmarshal(file, &facilities)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	r.facilities = facilities
	return nil
}

func (r *FacilityRegistry) GetFacility(id string) (Facility, bool) {
	return lo.Find(r.facilities, func(i Facility) bool {
		return i.ID == id
	})
}

func (r *FacilityRegistry) GetFacilityByTypeAndTier(facilityType FacilityType, tier int) (Facility, bool) {
	return lo.Find(r.facilities, func(i Facility) bool {
		return i.Tier == tier && i.Type == facilityType
	})
}

func (r *FacilityRegistry) GetAllFacilities() []Facility {
	return r.facilities
}

func (r *FacilityRegistry) GetAllFacilitiesByType(facilityType FacilityType) []Facility {
	var facilities []Facility
	for _, f := range r.facilities {
		if f.Type == facilityType {
			facilities = append(facilities, f)
		}
	}

	return facilities
}
