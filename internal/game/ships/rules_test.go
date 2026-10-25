package ships

import (
	"testing"

	"github.com/google/uuid"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRenameShip(t *testing.T) {
	agentID := uuid.New()
	ship := model.Ship{AgentID: agentID, Name: "old-name"}

	renamedShip := RenameShip(ship, "new-name")
	assert.Equal(t, "new-name", renamedShip.Name)
}

func TestChangeActiveShip(t *testing.T) {
	agentID := uuid.New()
	t.Run("success", func(t *testing.T) {
		oldActiveShip := model.Ship{AgentID: agentID, Active: true}
		newActiveShip := model.Ship{AgentID: agentID}

		newActiveShip, oldActiveShipToSave, _ := ChangeActiveShip(&oldActiveShip, newActiveShip)
		assert.True(t, newActiveShip.Active)
		if assert.NotNil(t, oldActiveShipToSave) {
			assert.False(t, oldActiveShipToSave.Active)
		}
	})

	t.Run("success: no old ship", func(t *testing.T) {
		newActiveShip, oldActiveShipToSave, _ := ChangeActiveShip(nil, model.Ship{AgentID: agentID})
		assert.True(t, newActiveShip.Active)
		assert.Nil(t, oldActiveShipToSave)
	})

	t.Run("invalid ship location", func(t *testing.T) {
		oldActiveShip := model.Ship{AgentID: agentID, Active: true, Coords: model.ShipCoords{SystemX: 1, SystemY: 1, Location: model.ShipLocationWaypoint, LocationID: 0}}
		newActiveShip := model.Ship{AgentID: agentID, Coords: model.ShipCoords{SystemX: 4, SystemY: 20, Location: model.ShipLocationPlanet, LocationID: 2}}
		_, _, err := ChangeActiveShip(&oldActiveShip, newActiveShip)
		assert.Error(t, err)
	})

	t.Run("invalid ship state", func(t *testing.T) {
		oldActiveShip := model.Ship{AgentID: agentID, Active: true, Status: model.ShipStatusOrbit}
		newActiveShip := model.Ship{AgentID: agentID, Status: model.ShipStatusDocked}
		_, _, err := ChangeActiveShip(&oldActiveShip, newActiveShip)
		assert.Error(t, err)
	})
}

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
	var withCode errs.WithCode
	if assert.ErrorAs(t, err, &withCode) {
		assert.Equal(t, withCode.Code, errs.CodeShipAlreadyInThisState)
	}
}

func TestDockShip(t *testing.T) {
	type testCase struct {
		name        string
		ship        model.Ship
		system      worldgen.System
		expectedErr bool
		err         error
		errCode     errs.ErrorCode
	}

	tests := []testCase{
		{
			name: "Success: waypoint",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					Location:   model.ShipLocationWaypoint,
					LocationID: 0,
				},
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
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					Location:   model.ShipLocationPlanet,
					LocationID: 0,
				},
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
			err:         errs.ErrUnprocessableEntity,
			errCode:     errs.CodeShipAlreadyInThisState,
		},
		{
			name: "Cant Dock",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					Location:   model.ShipLocationWaypoint,
					LocationID: 0,
				},
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
			err:         errs.ErrUnprocessableEntity,
			errCode:     errs.CodeCannotDock,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, err := DockShip(test.ship, test.system)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode errs.WithCode
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
