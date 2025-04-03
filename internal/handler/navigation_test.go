package handler

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
	planets, err := stateObj.S.FindAllPlanets(&model.Planet{Name: "testPlanet1"})
	assert.NoError(t, err)
	assert.Len(t, planets, 1)
	planet := planets[0]
	body, err := json.Marshal(&schema.FlyToPlanetSchema{PlanetID: planet.ID, SpaceshipID: spaceship.ID})
	assert.NoError(t, err)

	if !spaceship.PlayerSitIn {
		err = stateObj.S.EnterUserSpaceship(*usr, spaceship.ID)
		assert.NoError(t, err)
	}

	tests := []test.HTTPTest{
		{
			Description:   "Success",
			Method:        http.MethodPost,
			Route:         "/navigation/planet",
			Body:          body,
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var res schema.OkResponseSchema
				err = json.Unmarshal(body, &res)
				assert.NoError(t, err)

				assert.True(t, res.Ok)
				assert.Equal(t, 1, res.CustomStatusCode)
			},
		},
		{
			Description:   "Invalid body",
			Method:        http.MethodPost,
			Route:         "/navigation/planet",
			ExpectedError: true,
			ExpectedCode:  http.StatusUnprocessableEntity,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", userJwtToken),
	})

	stateObj.S.ExitUserSpaceship(*usr, spaceship.ID)
}

func TestHyperJump(t *testing.T) {
	flying := false
	err := stateObj.S.SetFlightInfo(spaceship.ID, &model.FlightInfo{Flying: &flying})
	assert.NoError(t, err)

	system, err := stateObj.S.CreateSystem(schema.CreateSystemSchema{Name: "hyperjumpSystem"})
	assert.NoError(t, err)

	body, err := json.Marshal(&schema.HyperJumpSchema{SystemID: system.ID, SpaceshipID: spaceship.ID})
	assert.NoError(t, err)

	if !spaceship.PlayerSitIn {
		err = stateObj.S.EnterUserSpaceship(*usr, spaceship.ID)
		assert.NoError(t, err)
	}

	tests := []test.HTTPTest{
		{
			Description:   "Success",
			Method:        http.MethodPost,
			Route:         "/navigation/hyperjump",
			Body:          body,
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var res schema.OkResponseSchema
				err = json.Unmarshal(body, &res)
				assert.NoError(t, err)

				assert.True(t, res.Ok)
				assert.Equal(t, 1, res.CustomStatusCode)
			},
		},
		{
			Description:   "Invalid body",
			Method:        http.MethodPost,
			Route:         "/navigation/hyperjump",
			ExpectedError: true,
			ExpectedCode:  http.StatusUnprocessableEntity,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", userJwtToken),
	})

	stateObj.S.ExitUserSpaceship(*usr, spaceship.ID)
}

func TestCheckFlight(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "Invalid spaceship UUID",
			Route:         "/navigation/info?id=123",
			ExpectedError: true,
			ExpectedCode:  http.StatusBadRequest,
			Method:        http.MethodGet,
		},
		{
			Description:   "Success",
			Route:         fmt.Sprintf("/navigation/info?id=%s", spaceship.ID.String()),
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var b *schema.FlyInfoSchema
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.Contains(t, []string{"planet", "system"}, b.Destination)
			},
			Method: http.MethodGet,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", userJwtToken),
	})
}
