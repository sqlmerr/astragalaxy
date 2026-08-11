package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/samber/lo"
)

type ResourceWorldGen struct {
	Planet struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"planet"`
}

type Resource struct {
	ID       string           `json:"id"`
	Tags     []string         `json:"tags"`
	WorldGen ResourceWorldGen `json:"worldgen"`
}

type ResourceRegistry struct {
	resources []Resource
}

func NewResourceRegistry() *ResourceRegistry {
	return &ResourceRegistry{resources: nil}
}

func (r *ResourceRegistry) Load(cfg Config) error {
	file, err := os.ReadFile(cfg.ResourcesPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	var resources []Resource
	err = json.Unmarshal(file, &resources)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	r.resources = resources
	return nil
}

func (r *ResourceRegistry) GetResource(id string) (Resource, bool) {
	return lo.Find(r.resources, func(i Resource) bool {
		return i.ID == id
	})
}

func (r *ResourceRegistry) GetAllResources() []Resource {
	return r.resources
}

func (r *ResourceRegistry) GetAllResourcesByTag(tag string) []Resource {
	var resources []Resource
	for _, res := range r.resources {
		if slices.Contains(res.Tags, tag) {
			resources = append(resources, res)
		}
	}
	return resources
}
