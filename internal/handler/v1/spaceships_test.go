package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"astragalaxy/pkg/test"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMySpaceships(t *testing.T) {
	api := createAPI(t)

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

	executor := test.New(api)
	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String()},
	)
}

func TestGetSpaceshipByID(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

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
				assert.Equal(t, testSpaceship.AstralID.String(), b["astral_id"].(string))
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

	executor.TestHTTP(
		t, tests,
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
			"X-Astral-ID":   testAstral.ID.String()},
	)
}

func TestSpaceshipOperations(t *testing.T) {
	t.Run("Exit", ExitMySpaceship)
	t.Run("Enter", EnterMySpaceship)
}

func EnterMySpaceship(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

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
			BeforeRequest: func() {
				err := testStateObj.S.ExitAstralSpaceship(testAstral.ID, testSpaceship.ID)
				assert.NoError(t, err)
				// time.Sleep(1)
			},
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
	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func ExitMySpaceship(t *testing.T) {
	api := createAPI(t)
	err := testStateObj.S.EnterAstralSpaceship(testAstral.ID, testSpaceship.ID)
	if !errors.Is(err, util.ErrPlayerAlreadyInSpaceship) {
		assert.NoError(t, err)
	}

	falseBool := false
	err = testStateObj.S.SetFlightInfo(testSpaceship.ID, &model.FlightInfo{Flying: &falseBool})
	assert.NoError(t, err)

	url := fmt.Sprintf("/v1/spaceships/my/%s/exit", testSpaceship.ID.String())

	res := api.Post(url, fmt.Sprintf("X-Astral-ID: %s", testAstral.ID.String()), fmt.Sprintf("Authorization: %s", testUserJwtToken))
	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		var response schema.OkResponse
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		s, err := testStateObj.S.FindOneSpaceship(testSpaceship.ID)
		assert.NoError(t, err)

		a, err := testStateObj.S.FindOneAstral(testAstral.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, response)
		assert.Equal(t, true, response.Ok)
		assert.Equal(t, 1, response.CustomStatusCode)
		assert.Equal(t, false, s.PlayerSitIn)
		assert.Equal(t, false, a.InSpaceship)
	}

	err = testStateObj.S.ExitAstralSpaceship(testAstral.ID, testSpaceship.ID)
	assert.NoError(t, err)
}

func TestRenameMySpaceship(t *testing.T) {
	api := createAPI(t)
	url := "/v1/spaceships/my/rename"
	body := &schema.RenameSpaceship{SpaceshipID: testSpaceship.ID, Name: "testSpaceship"}

	res := api.Patch(url, fmt.Sprintf("X-Astral-ID: %s", testAstral.ID.String()), fmt.Sprintf("Authorization: %s", testUserJwtToken), body)

	if assert.Equal(t, http.StatusOK, res.Code) {
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
