package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Facility struct {
	ID             string             `json:"id"`
	Type           model.FacilityType `json:"type"`
	Tier           int                `json:"tier"`
	TimeMultiplier float64            `json:"time_multiplier"`
	CostMultiplier float64            `json:"cost_multiplier"`
}

func facilityToModel(f Facility) model.Facility {
	return model.Facility{
		ID:             f.ID,
		Type:           model.FacilityType(f.Type),
		Tier:           f.Tier,
		TimeMultiplier: f.TimeMultiplier,
		CostMultiplier: f.CostMultiplier,
	}
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

func (r *FacilityRegistry) GetFacility(id string) (model.Facility, bool) {
	facility, found := lo.Find(r.facilities, func(i Facility) bool {
		return i.ID == id
	})
	return facilityToModel(facility), found
}

func (r *FacilityRegistry) GetFacilityByTypeAndTier(facilityType model.FacilityType, tier int) (model.Facility, bool) {
	facility, found := lo.Find(r.facilities, func(i Facility) bool {
		return i.Tier == tier && i.Type == facilityType
	})
	return facilityToModel(facility), found
}

func (r *FacilityRegistry) GetAllFacilities() []model.Facility {
	return lo.Map(r.facilities, func(f Facility, _ int) model.Facility {
		return facilityToModel(f)
	})
}

func (r *FacilityRegistry) GetAllFacilitiesByType(facilityType model.FacilityType) []model.Facility {
	var facilities []model.Facility
	for _, f := range r.facilities {
		if f.Type == facilityType {
			facilities = append(facilities, facilityToModel(f))
		}
	}

	return facilities
}
