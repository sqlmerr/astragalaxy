package handler

import (
	"astragalaxy/internal/schemas"
	"astragalaxy/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSystem(t *testing.T) {
	body := &schemas.CreateSystemSchema{Name: "testSystem"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/systems", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("secret-token", utils.Config("SECRET_TOKEN"))

	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusCreated, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response schemas.SystemSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, "testSystem", response.Name)
	}
}

func TestGetSystemPlanets(t *testing.T) {
	system, err := stateObj.SystemService.FindOneByName("initial")
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

		var response []schemas.PlanetSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Len(t, response, 0)
	}
}
