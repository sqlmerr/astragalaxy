package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type Facility struct {
	ID        string             `json:"id"`
	Type      model.FacilityType `json:"type"`
	Tier      int                `json:"tier"`
	Modifiers any                `json:"-"`
}

type FacilityCraftingModifiers struct {
	TimeMultiplier float64 `json:"time_multiplier"`
	CostMultiplier float64 `json:"cost_multiplier"`
}

func facilityToModel(f Facility) model.Facility {
	var mods any
	if craftingModifiers, ok := f.Modifiers.(*FacilityCraftingModifiers); ok {
		mods = &model.FacilityCraftingModifiers{
			TimeMultiplier: craftingModifiers.TimeMultiplier,
			CostMultiplier: craftingModifiers.CostMultiplier,
		}
	}
	return model.Facility{
		ID:        f.ID,
		Type:      model.FacilityType(f.Type),
		Tier:      f.Tier,
		Modifiers: mods,
	}
}

func (f *Facility) UnmarshalJSON(data []byte) error {
	// temp
	type Alias Facility
	aux := &struct {
		RawModifiers json.RawMessage `json:"modifiers"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch f.Type {
	case model.FacilityPrinter, model.FacilitySmelter:
		var m FacilityCraftingModifiers
		if err := json.Unmarshal(aux.RawModifiers, &m); err != nil {
			return err
		}
		f.Modifiers = &m

	case model.FacilityWalletTerminal:

	default:
		return fmt.Errorf("unknown facility type: %s", f.Type)
	}

	return nil
}

type FacilityRegistry struct {
	facilities      map[string]Facility
	facilitiesSlice []Facility
}

func NewFacilityRegistry() *FacilityRegistry {
	return &FacilityRegistry{facilities: nil, facilitiesSlice: nil}
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

	r.facilities = lo.SliceToMap(facilities, func(i Facility) (string, Facility) {
		return i.ID, i
	})
	r.facilitiesSlice = facilities
	return nil
}

func (r *FacilityRegistry) GetFacility(id string) (model.Facility, bool) {
	facility, exists := r.facilities[id]
	return facilityToModel(facility), exists
}

func (r *FacilityRegistry) GetFacilityByTypeAndTier(facilityType model.FacilityType, tier int) (model.Facility, bool) {
	facility, found := lo.Find(r.facilitiesSlice, func(i Facility) bool {
		return i.Tier == tier && i.Type == facilityType
	})
	return facilityToModel(facility), found
}

func (r *FacilityRegistry) GetAllFacilities() []model.Facility {
	return lo.Map(r.facilitiesSlice, func(f Facility, _ int) model.Facility {
		return facilityToModel(f)
	})
}

func (r *FacilityRegistry) GetAllFacilitiesByType(facilityType model.FacilityType) []model.Facility {
	var facilities []model.Facility
	for _, f := range r.facilitiesSlice {
		if f.Type == facilityType {
			facilities = append(facilities, facilityToModel(f))
		}
	}

	return facilities
}
