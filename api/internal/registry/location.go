package registry

import (
	"astragalaxy/internal/util"
	"encoding/json"
	"io"
	"os"
)

type Location struct {
	Code        string `json:"code"`
	Emoji       string `json:"emoji"`
	Multiplayer bool   `json:"multiplayer"`
}

type LocationRegistry struct {
	locations []Location
}

func NewLocation() LocationRegistry {
	return LocationRegistry{}
}

func (r *LocationRegistry) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var locations []Location

	err = json.Unmarshal(byteValue, &locations)
	if err != nil {
		return err
	}

	r.locations = locations

	return nil
}

func (r *LocationRegistry) All() []Location {
	return r.locations
}

func (r *LocationRegistry) FindOne(code string) (*Location, error) {
	for _, l := range r.locations {
		if l.Code == code {
			return &l, nil
		}
	}

	return nil, util.ErrNotFound
}
