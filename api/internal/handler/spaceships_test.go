package handler

import (
	"astragalaxy/internal/schemas"
	"astragalaxy/pkg/test"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMySpaceships(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/spaceships/my", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userJwtToken))

	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response []schemas.SpaceshipSchema

		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response)
		assert.Len(t, response, 1)
		assert.Equal(t, response[0].Name, spaceship.Name)
	}
}

func TestGetSpaceshipByID(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "spaceship found",
			Route:         fmt.Sprintf("/spaceships/%s", spaceship.ID),
			ExpectedError: false,
			ExpectedCode:  200,
			ExpectedBodyKeys: map[string]interface{}{
				"id":      spaceship.ID.String(),
				"name":    spaceship.Name,
				"user_id": spaceship.UserID.String(),
			},
			Method: http.MethodGet,
		},
		{
			Description:   "invalid id",
			Route:         "/spaceships/123",
			ExpectedError: true,
			ExpectedCode:  400,
			Method:        http.MethodGet,
		},
	}

	executor.TestHTTP(
		t, tests,
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", userJwtToken)},
	)
}

func TestEnterMySpaceship(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "entered spaceship",
			Route:         fmt.Sprintf("/spaceships/my/%s/enter", spaceship.ID),
			ExpectedError: false,
			ExpectedCode:  200,
			ExpectedBodyKeys: map[string]interface{}{
				"ok":                 true,
				"custom_status_code": float64(1),
			},
			Method: http.MethodPost,
		},
		{
			Description:   "can't enter spaceship",
			Route:         fmt.Sprintf("/spaceships/my/%s/enter", spaceship.ID),
			ExpectedError: true,
			ExpectedCode:  400,
			Method:        http.MethodPost,
		},
		{
			Description:   "invalid id",
			Route:         "/spaceships/my/123/enter",
			ExpectedError: true,
			ExpectedCode:  400,
			Method:        http.MethodPost,
		},
	}
	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", userJwtToken)})

	//url := fmt.Sprintf("/spaceships/my/%s/enter", spaceship.ID.String())
	//
	//req := httptest.NewRequest(http.MethodPost, url, nil)
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userJwtToken))
	//
	//res, err := app.Test(req, -1)
	//assert.NoError(t, err)
	//if assert.Equal(t, http.StatusOK, res.StatusCode) {
	//	body, err := io.ReadAll(res.Body)
	//	assert.NoError(t, err)
	//	var response schemas.OkResponseSchema
	//	err = json.Unmarshal(body, &response)
	//	assert.NoError(t, err)
	//
	//	s, err := stateObj.SpaceshipService.FindOne(spaceship.ID)
	//	assert.NoError(t, err)
	//
	//	p, err := stateObj.UserService.FindOne(usr.ID)
	//	assert.NoError(t, err)
	//
	//	assert.NotEmpty(t, response)
	//	assert.Equal(t, true, response.Ok)
	//	assert.Equal(t, 1, response.CustomStatusCode)
	//	assert.Equal(t, true, s.PlayerSitIn)
	//	assert.Equal(t, true, p.InSpaceship)
	//}
}

func TestExitMySpaceship(t *testing.T) {
	url := fmt.Sprintf("/spaceships/my/%s/exit", spaceship.ID.String())

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userJwtToken))

	res, err := app.Test(req, -1)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response schemas.OkResponseSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		s, err := stateObj.SpaceshipService.FindOne(spaceship.ID)
		assert.NoError(t, err)

		p, err := stateObj.UserService.FindOne(usr.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, response)
		assert.Equal(t, true, response.Ok)
		assert.Equal(t, 1, response.CustomStatusCode)
		assert.Equal(t, false, s.PlayerSitIn)
		assert.Equal(t, false, p.InSpaceship)
	}
}

func TestRenameMySpaceship(t *testing.T) {
	url := "/spaceships/my/rename"
	body := &schemas.RenameSpaceshipSchema{SpaceshipID: spaceship.ID, Name: "testSpaceship"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userJwtToken))

	res, err := app.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response schemas.OkResponseSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		s, err := stateObj.SpaceshipService.FindOne(spaceship.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, response)
		assert.Equal(t, true, response.Ok)
		assert.Equal(t, 1, response.CustomStatusCode)
		assert.Equal(t, "testSpaceship", s.Name)
	}
}
