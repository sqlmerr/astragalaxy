package navigation

import (
	"testing"

	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNavigateWarp(t *testing.T) {
	type testCase struct {
		name        string
		ship        model.Ship
		newSystem   worldgen.System
		fuelT1      int
		fuelT2      int
		expectedErr bool
		err         error
		errCode     errs.ErrorCode
	}

	tests := []testCase{
		{
			name: "Success",
			ship: model.Ship{
				Coords: model.ShipCoords{
					SystemX: 10,
					SystemY: 10,
				},
				Status: model.ShipStatusOrbit,
			},
			newSystem: worldgen.System{
				X: 10,
				Y: 20,
			},
			fuelT1:      4,
			expectedErr: false,
		},
		{
			name: "Success with tier 2",
			ship: model.Ship{
				Coords: model.ShipCoords{
					SystemX: 0,
					SystemY: 0,
				},
				Status: model.ShipStatusOrbit,
			},
			newSystem: worldgen.System{
				X: 18,
				Y: 0,
			},
			fuelT2:      2,
			expectedErr: false,
		},
		{
			name: "Not enough fuel",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					SystemX: 10,
					SystemY: 10,
				},
			},
			newSystem: worldgen.System{
				X: 10,
				Y: 20,
			},
			expectedErr: true,
			err:         errs.ErrUnprocessableEntity,
			errCode:     errs.CodeNotEnoughResources,
		},
		{
			name: "Invalid ship state",
			ship: model.Ship{
				Status: model.ShipStatusDocked,
			},
			newSystem:   worldgen.System{},
			expectedErr: true,
			err:         errs.ErrUnprocessableEntity,
			errCode:     errs.CodeInvalidShipState,
		},
		{
			name: "Invalid warp path",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					SystemX: -10,
					SystemY: -10,
				},
			},
			newSystem: worldgen.System{
				X: 10,
				Y: 10,
			},
			fuelT2:      12,
			expectedErr: true,
			err:         errs.ErrInvalidArgument,
			errCode:     errs.CodeInvalidWarpPath,
		},
		{
			name: "Already at destination",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					SystemX: 10,
					SystemY: 10,
				},
			},
			newSystem: worldgen.System{
				X: 10,
				Y: 10,
			},
			expectedErr: true,
			err:         errs.ErrNotModified,
			errCode:     errs.CodeAlreadyAtDestination,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, _, _, err := NavigateWarp(test.ship, test.newSystem, test.fuelT1, test.fuelT2)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode errs.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.newSystem.X, ship.Coords.SystemX)
				assert.Equal(t, test.newSystem.Y, ship.Coords.SystemY)
				assert.True(t, cd > 0)
			}
		})
	}
}

func TestNavigatePlanet(t *testing.T) {
	type testCase struct {
		name        string
		ship        model.Ship
		system      worldgen.System
		orbitIndex  int
		expectedErr bool
		err         error
		errCode     errs.ErrorCode
	}

	tests := []testCase{
		{
			name: "Success",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
			},
			system: worldgen.System{
				Planets: []worldgen.Planet{
					{
						Name:  "Planet-0",
						Orbit: 0,
					},
					{
						Name:  "Planet-1",
						Orbit: 1,
					},
				},
			},
			orbitIndex:  1,
			expectedErr: false,
		},
		{
			name: "Invalid ship state",
			ship: model.Ship{
				Status: model.ShipStatusDocked,
			},
			system:      worldgen.System{},
			orbitIndex:  1,
			expectedErr: true,
			err:         errs.ErrUnprocessableEntity,
			errCode:     errs.CodeInvalidShipState,
		},
		{
			name: "Invalid coordinates",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
			},
			system:      worldgen.System{},
			orbitIndex:  1,
			expectedErr: true,
			err:         errs.ErrNotFound,
			errCode:     errs.CodeInvalidCoordinates,
		},
		{
			name: "Already at destination",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					Location:   model.ShipLocationPlanet,
					LocationID: 1,
				},
			},
			system:      worldgen.System{},
			orbitIndex:  1,
			expectedErr: true,
			err:         errs.ErrNotModified,
			errCode:     errs.CodeAlreadyAtDestination,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, err := NavigatePlanet(test.ship, test.system, test.orbitIndex)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode errs.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, model.ShipLocationPlanet, ship.Coords.Location)
				assert.Equal(t, test.orbitIndex, ship.Coords.LocationID)
				assert.True(t, cd > 0)
			}
		})
	}
}

func TestNavigateWaypoint(t *testing.T) {
	type testCase struct {
		name        string
		ship        model.Ship
		system      worldgen.System
		waypointID  int
		expectedErr bool
		err         error
		errCode     errs.ErrorCode
	}

	tests := []testCase{
		{
			name: "Success",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
			},
			system: worldgen.System{
				Waypoints: []worldgen.Waypoint{
					{
						ID:   0,
						Type: worldgen.WaypointStation,
					},
				},
			},
			waypointID:  0,
			expectedErr: false,
		},
		{
			name: "Invalid ship state",
			ship: model.Ship{
				Status: model.ShipStatusDocked,
			},
			system:      worldgen.System{},
			waypointID:  0,
			expectedErr: true,
			err:         errs.ErrUnprocessableEntity,
			errCode:     errs.CodeInvalidShipState,
		},
		{
			name: "Invalid coordinates",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
			},
			system:      worldgen.System{},
			waypointID:  0,
			expectedErr: true,
			err:         errs.ErrNotFound,
			errCode:     errs.CodeInvalidCoordinates,
		},
		{
			name: "Already at destination",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
				Coords: model.ShipCoords{
					Location:   model.ShipLocationWaypoint,
					LocationID: 0,
				},
			},
			system:      worldgen.System{},
			waypointID:  0,
			expectedErr: true,
			err:         errs.ErrNotModified,
			errCode:     errs.CodeAlreadyAtDestination,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, err := NavigateWaypoint(test.ship, test.system, test.waypointID)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode errs.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, model.ShipLocationWaypoint, ship.Coords.Location)
				assert.Equal(t, test.waypointID, ship.Coords.LocationID)
				assert.True(t, cd > 0)
			}
		})
	}
}
