package v1

import (
	"astragalaxy/internal/registry"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItemByCode(t *testing.T) {
	api := createAPI(t)
	resNotFound := api.Get("/v1/registry/items/notfound")
	assert.Equal(t, http.StatusNotFound, resNotFound.Code)

	code := "test"
	url := fmt.Sprintf("/v1/registry/items/%s", code)
	res := api.Get(url)

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, _ := io.ReadAll(res.Body)
		var item registry.RItem
		err := json.Unmarshal(body, &item)
		assert.NoError(t, err)

		assert.Equal(t, code, item.Code)
		assert.Equal(t, registry.ITEM_RARITY_IMMORTAL, item.Rarity)
	}
}

func TestGetItems(t *testing.T) {
	api := createAPI(t)

	url := "/v1/registry/items"
	res := api.Get(url)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetLocations(t *testing.T) {
	api := createAPI(t)

	url := "/v1/registry/locations"
	res := api.Get(url)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetLocationByCode(t *testing.T) {
	api := createAPI(t)
	resNotFound := api.Get("/v1/registry/locations/notfound")
	assert.Equal(t, http.StatusNotFound, resNotFound.Code)

	code := "space_station"
	url := fmt.Sprintf("/v1/registry/locations/%s", code)
	res := api.Get(url)

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, _ := io.ReadAll(res.Body)
		var location registry.Location
		err := json.Unmarshal(body, &location)
		assert.NoError(t, err)

		assert.Equal(t, code, location.Code)
	}
}
