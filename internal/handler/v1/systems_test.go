package v1

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
	body := &schema.CreateSystem{Name: "testSystem"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/systems", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("secret-token", testStateObj.Config.SecretToken)

	res, err := testApp.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusCreated, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response schema.System
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, "testSystem", response.Name)
	}
}

func TestGetSystemPlanets(t *testing.T) {
	system, err := testStateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/systems/%s/planets", system.ID), nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", testUserJwtToken))

	res, err := testApp.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
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
	system, err := testStateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/systems/%s", system.ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", testUserJwtToken))
	res, err := testApp.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
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
	system, err := testStateObj.S.FindOneSystemByName("initial")
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/v1/systems", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", testUserJwtToken))
	res, err := testApp.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
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
