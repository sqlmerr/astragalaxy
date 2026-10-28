package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type ShipClass struct {
	ID model.ShipClassID `json:"id"`

	MaxModuleSlots      int     `json:"max_module_slots"`
	MaxWarpDistance     int     `json:"max_warp_distance"`
	MiningLimitModifier float64 `json:"mining_limit_modifier"`
}

func shipClassToModel(m ShipClass) model.ShipClass {
	return model.ShipClass{
		ID:                  model.ShipClassID(m.ID),
		MaxModuleSlots:      m.MaxModuleSlots,
		MaxWarpDistance:     m.MaxWarpDistance,
		MiningLimitModifier: m.MiningLimitModifier,
	}
}

type ShipClassRegistry struct {
	classes []ShipClass
}

func NewShipClassRegistry() *ShipClassRegistry {
	return &ShipClassRegistry{classes: nil}
}

func (r *ShipClassRegistry) Load(cfg Config) error {
	file, err := os.ReadFile(cfg.ShipClassesPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	var classes []ShipClass
	err = json.Unmarshal(file, &classes)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	r.classes = classes
	return nil
}

func (r *ShipClassRegistry) GetShipClass(id model.ShipClassID) (model.ShipClass, bool) {
	class, found := lo.Find(r.classes, func(i ShipClass) bool {
		return i.ID == id
	})
	return shipClassToModel(class), found
}

func (r *ShipClassRegistry) GetShipClasses() []model.ShipClass {
	return lo.Map(r.classes, func(class ShipClass, _ int) model.ShipClass {
		return shipClassToModel(class)
	})
}
