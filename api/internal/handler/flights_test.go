package handler

import (
	"astragalaxy/internal/models"
	"astragalaxy/internal/schemas"
	"astragalaxy/pkg/test"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFlightToPlanet(t *testing.T) {
	planets, err := stateObj.PlanetService.FindAll(&models.Planet{Name: "testPlanet1"})
	assert.NoError(t, err)
	assert.Len(t, planets, 1)
	planet := planets[0]
	body, err := json.Marshal(&schemas.FlyToPlanetSchema{PlanetID: planet.ID, SpaceshipID: spaceship.ID})
	assert.NoError(t, err)

	if !spaceship.PlayerSitIn {
		err = stateObj.UserService.EnterSpaceship(*usr, spaceship.ID)
		assert.NoError(t, err)
	}

	tests := []test.HTTPTest{
		{
			Description:   "Success",
			Method:        http.MethodPost,
			Route:         "/flights/planet",
			Body:          body,
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var res schemas.OkResponseSchema
				err = json.Unmarshal(body, &res)
				assert.NoError(t, err)

				assert.True(t, res.Ok)
				assert.Equal(t, 1, res.CustomStatusCode)
			},
		},
		{
			Description:   "Invalid body",
			Method:        http.MethodPost,
			Route:         "/flights/planet",
			ExpectedError: true,
			ExpectedCode:  http.StatusUnprocessableEntity,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", userJwtToken),
	})

	stateObj.UserService.ExitSpaceship(*usr, spaceship.ID)
}

func TestCheckFlight(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "Invalid spaceship UUID",
			Route:         "/flights/info?id=123",
			ExpectedError: true,
			ExpectedCode:  http.StatusBadRequest,
			Method:        http.MethodGet,
		},
		{
			Description:   "Success",
			Route:         fmt.Sprintf("/flights/info?id=%s", spaceship.ID.String()),
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(body []byte) {
				var b *schemas.FlyInfoSchema
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
			},
			Method: http.MethodGet,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", userJwtToken),
	})
}
