package handler

import (
	"astragalaxy/internal/schema"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSystem(t *testing.T) {
	body := &schema.CreateSystemSchema{Name: "testSystem"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/systems", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("secret-token", stateObj.Config.SecretToken)

	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusCreated, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response schema.SystemSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, "testSystem", response.Name)
	}
}

func TestGetSystemPlanets(t *testing.T) {
	system, err := stateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/systems/%s/planets", system.ID.String()), nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userJwtToken))

	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response []schema.PlanetSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Len(t, response, 1)
		planet := response[0]
		assert.Equal(t, "testPlanet1", planet.Name)
		assert.Equal(t, "TOXINS", planet.Threat)
	}
}

func TestGetSystemByID(t *testing.T) {
	system, err := stateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/systems/%s", system.ID.String()), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userJwtToken))
	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response schema.SystemSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)
		assert.Equal(t, "initial", response.Name)
		assert.Equal(t, system.ID, response.ID)
	}
}

func TestGetAllSystems(t *testing.T) {
	system, err := stateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/systems", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", userJwtToken))
	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response []schema.SystemSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, system.ID, response[0].ID)
		assert.Equal(t, system.Name, response[0].Name)
	}
}
