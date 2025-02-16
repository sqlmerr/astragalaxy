package handler

import (
	"astragalaxy/internal/registry"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItemByCode(t *testing.T) {
	resNotFound, err := app.Test(httptest.NewRequest(http.MethodGet, "/registry/items/notfound", nil), -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resNotFound.StatusCode)

	code := "teleporter"
	url := fmt.Sprintf("/registry/items/%s", code)
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	res, err := app.Test(request, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, _ := io.ReadAll(res.Body)
		var item registry.Item
		err := json.Unmarshal(body, &item)
		assert.NoError(t, err)

		assert.Equal(t, code, item.Code)
	}
}

func TestGetItems(t *testing.T) {
	url := "/registry/items"
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	res, err := app.Test(request, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
