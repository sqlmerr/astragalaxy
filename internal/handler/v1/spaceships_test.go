package v1

import (
	"astragalaxy/internal/schema"
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
	tests := []test.HTTPTest{
		{
			Description:   "Should return 200 OK response",
			Method:        http.MethodGet,
			ExpectedCode:  http.StatusOK,
			ExpectedError: false,
			Route:         "/v1/spaceships/my",
			BodyValidator: func(body []byte) {
				var res schema.DataGenericResponse[[]schema.Spaceship]
				err := json.Unmarshal(body, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res.Data)

				if assert.Len(t, res.Data, 1) {
					assert.Equal(t, res.Data[0].Name, testSpaceship.Name)
				}
			},
		},
	}

	testExecutor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken)},
	)
}

func TestGetSpaceshipByID(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "testSpaceship found",
			Route:         fmt.Sprintf("/v1/spaceships/%s", testSpaceship.ID),
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b map[string]interface{}
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.Equal(t, testSpaceship.ID.String(), b["id"].(string))
				assert.Equal(t, testSpaceship.Name, b["name"].(string))
				assert.Equal(t, testSpaceship.UserID.String(), b["user_id"].(string))
			},
			Method: http.MethodGet,
		},
		{
			Description:   "invalid id",
			Route:         "/v1/spaceships/123",
			ExpectedError: true,
			ExpectedCode:  400,
			Method:        http.MethodGet,
		},
	}

	testExecutor.TestHTTP(
		t, tests,
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken)},
	)
}

func TestEnterMySpaceship(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "entered testSpaceship",
			Route:         fmt.Sprintf("/v1/spaceships/my/%s/enter", testSpaceship.ID),
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b map[string]interface{}
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.Equal(t, true, b["ok"])
				assert.Equal(t, float64(1), b["custom_status_code"])
			},
			Method: http.MethodPost,
		},
		{
			Description:   "can't enter testSpaceship",
			Route:         fmt.Sprintf("/v1/spaceships/my/%s/enter", testSpaceship.ID),
			ExpectedError: true,
			ExpectedCode:  400,
			Method:        http.MethodPost,
		},
		{
			Description:   "invalid id",
			Route:         "/v1/spaceships/my/123/enter",
			ExpectedError: true,
			ExpectedCode:  400,
			Method:        http.MethodPost,
		},
	}
	testExecutor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken)})
}

func TestExitMySpaceship(t *testing.T) {
	url := fmt.Sprintf("/v1/spaceships/my/%s/exit", testSpaceship.ID.String())

	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testUserJwtToken))

	res, err := testApp.Test(req, -1)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response schema.OkResponse
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		s, err := testStateObj.S.FindOneSpaceship(testSpaceship.ID)
		assert.NoError(t, err)

		p, err := testStateObj.S.FindOneUser(testUser.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, response)
		assert.Equal(t, true, response.Ok)
		assert.Equal(t, 1, response.CustomStatusCode)
		assert.Equal(t, false, s.PlayerSitIn)
		assert.Equal(t, false, p.InSpaceship)
	}
}

func TestRenameMySpaceship(t *testing.T) {
	url := "/v1/spaceships/my/rename"
	body := &schema.RenameSpaceship{SpaceshipID: testSpaceship.ID, Name: "testSpaceship"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testUserJwtToken))

	res, err := testApp.Test(req, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response schema.OkResponse
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		s, err := testStateObj.S.FindOneSpaceship(testSpaceship.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, response)
		assert.Equal(t, true, response.Ok)
		assert.Equal(t, 1, response.CustomStatusCode)
		assert.Equal(t, "testSpaceship", s.Name)
	}
}
