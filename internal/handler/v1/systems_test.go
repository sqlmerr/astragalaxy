package v1

import (
	"astragalaxy/internal/schema"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSystem(t *testing.T) {
	api := createAPI(t)
	body := &schema.CreateSystem{Name: "testSystem"}
	res := api.Post("/v1/systems/", fmt.Sprintf("secret-token: %s", testStateObj.Config.SecretToken), body)

	if assert.Equal(t, http.StatusCreated, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response schema.System
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, "testSystem", response.Name)
	}
}

func TestGetSystemPlanets(t *testing.T) {
	api := createAPI(t)
	system, err := testStateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)
	url := fmt.Sprintf("/v1/systems/%s/planets", system.ID)
	assert.NoError(t, err)

	res := api.Get(url, fmt.Sprintf("Authorization: %s", fmt.Sprintf("Bearer %v", testUserJwtToken)))

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response schema.DataGenericResponse[[]schema.Planet]
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		if assert.Len(t, response.Data, 1) {
			planet := response.Data[0]
			assert.Equal(t, "testPlanet1", planet.Name)
			assert.Equal(t, "TOXINS", planet.Threat)
		}
	}
}

func TestGetSystemByID(t *testing.T) {
	api := createAPI(t)
	system, err := testStateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)
	url := fmt.Sprintf("/v1/systems/%s", system.ID)

	res := api.Get(url, fmt.Sprintf("Authorization: %s", fmt.Sprintf("Bearer %v", testUserJwtToken)))

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response schema.System
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, "initial", response.Name)
		assert.Equal(t, system.ID, response.ID)
	}
}

func TestGetAllSystems(t *testing.T) {
	api := createAPI(t)
	system, err := testStateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)
	res := api.Get("/v1/systems/", fmt.Sprintf("Authorization: %s", fmt.Sprintf("Bearer %v", testUserJwtToken)))

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var r schema.DataGenericResponse[[]schema.System]
		err = json.Unmarshal(body, &r)
		assert.NoError(t, err)

		if assert.GreaterOrEqual(t, len(r.Data), 1) {
			assert.Equal(t, system.ID, r.Data[0].ID)
			assert.Equal(t, system.Name, r.Data[0].Name)
		}
	}
}
