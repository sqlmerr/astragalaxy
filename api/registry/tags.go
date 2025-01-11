package registry

import (
	"astragalaxy/utils"
	"encoding/json"
	"io"
	"os"
	"slices"
)

type Tag struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

type TagRegistry struct {
	tags         []Tag
	itemRegistry ItemRegistry
}

func NewTag(itemRegistry ItemRegistry) TagRegistry {
	return TagRegistry{
		itemRegistry: itemRegistry,
	}
}

func (r *TagRegistry) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var tags []Tag

	err = json.Unmarshal(byteValue, &tags)
	if err != nil {
		return err
	}

	r.tags = tags

	return nil
}

func (r *TagRegistry) HasItemTag(itemCode string, tagName string) bool {
	tag, err := r.FindOne(tagName)
	if err != nil {
		return false
	}
	if !slices.Contains(tag.Items, itemCode) {
		return false
	}
	return true
}

func (r *TagRegistry) FindOne(tag string) (*Tag, error) {
	for _, t := range r.tags {
		if t.Name == tag {
			return &t, nil
		}
	}

	return nil, utils.ErrItemNotFound
}
