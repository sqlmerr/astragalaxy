package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/internal/util"
	"astragalaxy/pkg/test"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartExploration(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	types := []string{"gathering", "mining", "structures"}
	planetExploration := types[rand.IntN(len(types)-1)]

	tests := []test.HTTPTest{
		{
			Description:   "Invalid Type",
			Method:        http.MethodPost,
			Route:         "/v1/explore/teststst123",
			ExpectedError: true,
			ExpectedCode:  http.StatusUnprocessableEntity,
		},
		{
			Description:   "Success at planet",
			Method:        http.MethodPost,
			Route:         fmt.Sprintf("/v1/explore/%s", planetExploration),
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(bytes []byte) {
				var b schema.OkResponse
				err := json.Unmarshal(bytes, &b)
				assert.NoError(t, err)

				assert.Equal(t, true, b.Ok)

				info, err := testStateObj.S.GetExplorationInfoOrCreate(testAstral.ID)
				if assert.NoError(t, err) {
					assert.Equal(t, true, info.Status)
					assert.Equal(t, planetExploration, string(info.Type))
				}
			},
			AfterRequest: func() {
				info, err := testStateObj.S.GetExplorationInfoOrCreate(testAstral.ID)
				assert.NoError(t, err)

				err = testStateObj.S.SetExplorationInfo(info.ID, map[string]any{"exploring": false, "type": nil, "started_at": 0, "required_time": 0})
				assert.NoError(t, err)
			},
			BeforeRequest: func() {
				err := testStateObj.S.EnterAstralSpaceship(testAstral.ID, testSpaceship.ID)
				if !errors.Is(err, util.ErrPlayerAlreadyInSpaceship) {
					assert.NoError(t, err)
				}

				err = testStateObj.S.UpdateSpaceship(testSpaceship.ID, schema.UpdateSpaceship{PlanetID: testPlanet.ID, Location: "planet"})
				assert.NoError(t, err)

				err = testStateObj.S.UpdateAstral(testAstral.ID, schema.UpdateAstral{Location: "planet"})
				assert.NoError(t, err)

				err = testStateObj.S.ExitAstralSpaceship(testAstral.ID, testSpaceship.ID)
				assert.NoError(t, err)
			},
		},
		{
			Description:   "Success at space",
			Method:        http.MethodPost,
			Route:         "/v1/explore/asteroids",
			ExpectedError: false,
			ExpectedCode:  http.StatusOK,
			BodyValidator: func(bytes []byte) {
				var b schema.OkResponse
				err := json.Unmarshal(bytes, &b)
				assert.NoError(t, err)

				assert.Equal(t, true, b.Ok)

				info, err := testStateObj.S.GetExplorationInfoOrCreate(testAstral.ID)
				if assert.NoError(t, err) {
					assert.Equal(t, true, info.Status)
					assert.Equal(t, "asteroids", string(info.Type))
				}
			},
			AfterRequest: func() {
				info, err := testStateObj.S.GetExplorationInfoOrCreate(testAstral.ID)
				assert.NoError(t, err)

				err = testStateObj.S.SetExplorationInfo(info.ID, map[string]any{"exploring": false, "type": nil, "started_at": 0, "required_time": 0})
				assert.NoError(t, err)

				// err = testStateObj.S.NavigateLocation(testSpaceship.ID, "open_space")
				// assert.NoError(t, err)
			},
			BeforeRequest: func() {
				err := testStateObj.S.EnterAstralSpaceship(testAstral.ID, testSpaceship.ID)
				assert.NoError(t, err)

				err = testStateObj.S.NavigateLocation(testSpaceship.ID, "open_space")
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
