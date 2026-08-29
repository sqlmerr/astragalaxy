package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type ResourceWorldGenParams struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Resource struct {
	ID       string                            `json:"id"`
	Tags     []string                          `json:"tags"`
	WorldGen map[string]ResourceWorldGenParams `json:"worldgen"`
}

func resourceToModel(r Resource) model.ResourceData {
	return model.ResourceData{
		ID:   r.ID,
		Tags: r.Tags,
		WorldGen: lo.MapValues(r.WorldGen, func(v ResourceWorldGenParams, k string) model.ResourceDataWorldGenParams {
			return model.ResourceDataWorldGenParams{
				Min: v.Min,
				Max: v.Max,
			}
		}),
	}
}

func (r *Resource) Normalize() {
	if r.Tags == nil {
		r.Tags = []string{}
	}
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
	var normalizedResources []Resource
	for _, r := range resources {
		r.Normalize()
		normalizedResources = append(normalizedResources, r)
	}

	r.resources = normalizedResources
	return nil
}

func (r *ResourceRegistry) GetResource(id string) (model.ResourceData, bool) {
	resource, found := lo.Find(r.resources, func(i Resource) bool {
		return i.ID == id
	})
	return resourceToModel(resource), found
}

func (r *ResourceRegistry) GetAllResources() []model.ResourceData {
	return lo.Map(r.resources, func(resource Resource, _ int) model.ResourceData {
		return resourceToModel(resource)
	})
}

func (r *ResourceRegistry) GetAllResourcesByTag(tag string) []model.ResourceData {
	var resources []model.ResourceData
	for _, res := range r.resources {
		if slices.Contains(res.Tags, tag) {
			resources = append(resources, resourceToModel(res))
		}
	}
	return resources
}

func (r *ResourceRegistry) GetAllResourcesByWorldgenParams(paramName string) []model.ResourceData {
	var resources []model.ResourceData
	for _, res := range r.resources {
		_, exists := res.WorldGen[paramName]
		if exists {
			resources = append(resources, resourceToModel(res))
		}
	}
	return resources
}
