package registry

import (
	"astragalaxy/internal/util"
	"encoding/json"
	"io"
	"os"

	"github.com/samber/lo"
)

type RResource struct {
	ID     string   `json:"id"`
	Groups []string `json:"groups"`
}

type ResourceRegistry struct {
	resources []RResource
}

func NewResource() ResourceRegistry {
	return ResourceRegistry{}
}

func (r *ResourceRegistry) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var resources []RResource

	err = json.Unmarshal(byteValue, &resources)
	if err != nil {
		return err
	}

	r.resources = resources

	return nil
}

func (r *ResourceRegistry) All() []RResource {
	return r.resources
}

func (r *ResourceRegistry) FindOne(id string) (*RResource, error) {
	for _, i := range r.resources {
		if i.ID == id {
			return &i, nil
		}
	}

	return nil, util.ErrNotFound
}

func (r *ResourceRegistry) FindAllByGroup(group string) []RResource {
	resources := make([]RResource, 0)
	for _, i := range r.resources {
		if lo.Contains(i.Groups, group) {
			resources = append(resources, i)
		}
	}

	return resources
}
