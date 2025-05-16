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
	api := createAPI(t)
	planets, err := testStateObj.S.FindAllPlanets(&model.Planet{Name: "testPlanet1"})
	assert.NoError(t, err)
	assert.Len(t, planets, 1)
	planet := planets[0]
	body := &schema.FlyToPlanet{PlanetID: planet.ID, SpaceshipID: testSpaceship.ID}

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
			ExpectedCode:  http.StatusBadRequest,
		},
	}

	executor := test.New(api)
	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String(),
	})

	testStateObj.S.ExitAstralSpaceship(*testAstral, testSpaceship.ID)
}

func TestHyperJump(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	flying := false
	err := testStateObj.S.SetFlightInfo(testSpaceship.ID, &model.FlightInfo{Flying: &flying})
	assert.NoError(t, err)

	system, err := testStateObj.S.CreateSystem(schema.CreateSystem{Name: "hyperjumpSystem", Connections: []string{testSpaceship.SystemID}})
	assert.NoError(t, err)

	body := schema.HyperJump{Path: fmt.Sprintf("%s->%s", testSpaceship.SystemID, system.ID), SpaceshipID: testSpaceship.ID}

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
			BeforeRequest: func() {
				
			},
		},
		{
			Description:   "Invalid body",
			Method:        http.MethodPost,
			Route:         "/v1/navigation/hyperjump",
			ExpectedError: true,
			ExpectedCode:  http.StatusBadRequest,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String(),
	})

	testStateObj.S.ExitAstralSpaceship(*testAstral, testSpaceship.ID)
	err = testStateObj.S.DeleteSystem(system.ID)
	assert.NoError(t, err)
}

func TestNavigateToLocation(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	tests := []test.HTTPTest{
		{
			Description:   "Invalid body",
			Method:        http.MethodPost,
			Route:         "/v1/navigation/location",
			ExpectedError: true,
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Description:   "Success",
			Method:        http.MethodPost,
			Route:         "/v1/navigation/location",
			Body:          &schema.NavigateToLocation{SpaceshipID: testSpaceship.ID, Location: "test"},
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var res schema.OkResponse
				err := json.Unmarshal(body, &res)
				assert.NoError(t, err)

				assert.True(t, res.Ok)
				assert.Equal(t, 1, res.CustomStatusCode)

				sys, err := testStateObj.S.FindOneSystem(testSpaceship.SystemID)
				assert.NoError(t, err)
				fmt.Println("    ar aradasd: ", sys.Locations)

				spcship, err := testStateObj.S.FindOneSpaceship(testSpaceship.ID)
				if assert.NoError(t, err) {
					assert.Equal(t, "test", spcship.Location)
				}
			},
			BeforeRequest: func() {
				err := testStateObj.S.UpdateSystem(testSpaceship.SystemID, schema.UpdateSystem{Locations: []string{"space_station", "test"}})
				assert.NoError(t, err)
			},
			AfterRequest: func() {
				err := testStateObj.S.UpdateSystem(testSpaceship.SystemID, schema.UpdateSystem{Locations: []string{"space_station"}})
				assert.NoError(t, err)
			},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String(),
	})
}

func TestCheckFlight(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

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
				if assert.NoError(t, err) {
					assert.Contains(t, []string{"planet", "system"}, b.Destination)
				}

			},
			Method: http.MethodGet,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken),
		"X-Astral-ID":   testAstral.ID.String(),
	})
}
