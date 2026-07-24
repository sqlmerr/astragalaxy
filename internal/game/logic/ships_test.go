package logic

import (
	"testing"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/stretchr/testify/assert"
)

func TestOrbitShip(t *testing.T) {
	ship1 := model.Ship{
		Status: model.ShipStatusDocked,
	}
	ship1, cd, err := OrbitShip(ship1)
	assert.NoError(t, err)
	assert.True(t, cd > 0)
	assert.Equal(t, ship1.Status, model.ShipStatusOrbit)

	ship2 := model.Ship{
		Status: model.ShipStatusOrbit,
	}
	_, _, err = OrbitShip(ship2)
	assert.Error(t, err)
	var withCode core_errors.WithCode
	if assert.ErrorAs(t, err, &withCode) {
		assert.Equal(t, withCode.Code, core_errors.CodeShipAlreadyInThisState)
	}
}

func TestDockShip(t *testing.T) {
	type testCase struct {
		name        string
		ship        model.Ship
		system      worldgen.System
		expectedErr bool
		err         error
		errCode     core_errors.ErrorCode
	}

	tests := []testCase{
		{
			name: "Success: waypoint",
			ship: model.Ship{
				Status:     model.ShipStatusOrbit,
				Location:   model.ShipLocationWaypoint,
				LocationID: 0,
			},
			system: worldgen.System{
				Name: "System-1",
				Waypoints: []worldgen.Waypoint{
					{
						ID:       0,
						Type:     worldgen.WaypointStation,
						Dockable: true,
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "Success: planet",
			ship: model.Ship{
				Status:     model.ShipStatusOrbit,
				Location:   model.ShipLocationPlanet,
				LocationID: 0,
			},
			system: worldgen.System{
				Name: "System-2",
				Planets: []worldgen.Planet{
					{
						Name:  "Planet",
						Type:  worldgen.PlanetTerra,
						Orbit: 0,
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "Ship Already Docked",
			ship: model.Ship{
				Status: model.ShipStatusDocked,
			},
			system:      worldgen.System{},
			expectedErr: true,
			err:         core_errors.ErrUnprocessableEntity,
			errCode:     core_errors.CodeShipAlreadyInThisState,
		},
		{
			name: "Cant Dock",
			ship: model.Ship{
				Status:     model.ShipStatusOrbit,
				Location:   model.ShipLocationWaypoint,
				LocationID: 0,
			},
			system: worldgen.System{
				Name: "System-3",
				Waypoints: []worldgen.Waypoint{
					{
						ID:       0,
						Type:     worldgen.WaypointAsteroid,
						Dockable: false,
					},
				},
			},
			expectedErr: true,
			err:         core_errors.ErrUnprocessableEntity,
			errCode:     core_errors.CodeCannotDock,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, err := DockShip(test.ship, test.system)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.True(t, cd > 0)
				assert.Equal(t, model.ShipStatusDocked, ship.Status)
			}
		})
	}
}
