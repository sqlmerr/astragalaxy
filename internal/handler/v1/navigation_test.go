package v1

import (
	"astragalaxy/internal/model"
	"astragalaxy/internal/schema"
	"astragalaxy/pkg/test"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlightToPlanet(t *testing.T) {
	planets, err := testStateObj.S.FindAllPlanets(&model.Planet{Name: "testPlanet1"})
	assert.NoError(t, err)
	assert.Len(t, planets, 1)
	planet := planets[0]
	body, err := json.Marshal(&schema.FlyToPlanet{PlanetID: planet.ID, SpaceshipID: testSpaceship.ID})
	assert.NoError(t, err)

	if !testSpaceship.PlayerSitIn {
		err = testStateObj.S.EnterAstralSpaceship(*testAstral, testSpaceship.ID)
		assert.NoError(t, err)
	}

	tests := []test.HTTPTest{
		{
			Description:   "Success",
			Method:        http.MethodPost,
			Route:         "/v1/navigation/planet",
			Body:          body,
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var res schema.OkResponse
				err = json.Unmarshal(body, &res)
				assert.NoError(t, err)

				assert.True(t, res.Ok)
				assert.Equal(t, 1, res.CustomStatusCode)
			},
		},
		{
			Description:   "Invalid body",
			Method:        http.MethodPost,
			Route:         "/v1/navigation/planet",
			ExpectedError: true,
			ExpectedCode:  http.StatusUnprocessableEntity,
		},
	}

	testExecutor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String(),
	})

	testStateObj.S.ExitAstralSpaceship(*testAstral, testSpaceship.ID)
}

func TestHyperJump(t *testing.T) {
	flying := false
	err := testStateObj.S.SetFlightInfo(testSpaceship.ID, &model.FlightInfo{Flying: &flying})
	assert.NoError(t, err)

	system, err := testStateObj.S.CreateSystem(schema.CreateSystem{Name: "hyperjumpSystem", Connections: []string{testSpaceship.SystemID}})
	assert.NoError(t, err)

	body, err := json.Marshal(&schema.HyperJump{Path: fmt.Sprintf("%s->%s", testSpaceship.SystemID, system.ID), SpaceshipID: testSpaceship.ID})
	assert.NoError(t, err)

	if !testSpaceship.PlayerSitIn {
		err = testStateObj.S.EnterAstralSpaceship(*testAstral, testSpaceship.ID)
		assert.NoError(t, err)
	}

	tests := []test.HTTPTest{
		{
			Description:   "Success",
			Method:        http.MethodPost,
			Route:         "/v1/navigation/hyperjump",
			Body:          body,
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var res schema.OkResponse
				err = json.Unmarshal(body, &res)
				assert.NoError(t, err)

				assert.True(t, res.Ok)
				assert.Equal(t, 1, res.CustomStatusCode)
			},
		},
		{
			Description:   "Invalid body",
			Method:        http.MethodPost,
			Route:         "/v1/navigation/hyperjump",
			ExpectedError: true,
			ExpectedCode:  http.StatusUnprocessableEntity,
		},
	}

	testExecutor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String(),
	})

	testStateObj.S.ExitAstralSpaceship(*testAstral, testSpaceship.ID)
}

func TestCheckFlight(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "Invalid testSpaceship UUID",
			Route:         "/v1/navigation/info?id=123",
			ExpectedError: true,
			ExpectedCode:  http.StatusBadRequest,
			Method:        http.MethodGet,
		},
		{
			Description:   "Success",
			Route:         fmt.Sprintf("/v1/navigation/info?id=%s", testSpaceship.ID.String()),
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var b *schema.FlyInfo
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.Contains(t, []string{"planet", "system"}, b.Destination)
			},
			Method: http.MethodGet,
		},
	}

	testExecutor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String(),
	})
}
