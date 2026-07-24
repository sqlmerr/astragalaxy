package logic

import (
	"testing"

	"github.com/sqlmerr/astragalaxy/internal/data/model"
	core_errors "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/stretchr/testify/assert"
)

func TestNavigateWarp(t *testing.T) {
	type testCase struct {
		name        string
		ship        model.Ship
		newSystem   worldgen.System
		expectedErr bool
		err         error
		errCode     core_errors.ErrorCode
	}

	tests := []testCase{
		{
			name: "Success",
			ship: model.Ship{
				SystemX: 10,
				SystemY: 10,
				Status:  model.ShipStatusOrbit,
			},
			newSystem: worldgen.System{
				X: 10,
				Y: 20,
			},
			expectedErr: false,
		},
		{
			name: "Invalid ship state",
			ship: model.Ship{
				Status: model.ShipStatusDocked,
			},
			newSystem:   worldgen.System{},
			expectedErr: true,
			err:         core_errors.ErrUnprocessableEntity,
			errCode:     core_errors.CodeInvalidShipState,
		},
		{
			name: "Invalid warp path",
			ship: model.Ship{
				Status:  model.ShipStatusOrbit,
				SystemX: -10,
				SystemY: -10,
			},
			newSystem: worldgen.System{
				X: 10,
				Y: 10,
			},
			expectedErr: true,
			err:         core_errors.ErrInvalidArgument,
			errCode:     core_errors.CodeInvalidWarpPath,
		},
		{
			name: "Already at destination",
			ship: model.Ship{
				Status:  model.ShipStatusOrbit,
				SystemX: 10,
				SystemY: 10,
			},
			newSystem: worldgen.System{
				X: 10,
				Y: 10,
			},
			expectedErr: true,
			err:         core_errors.ErrNotModified,
			errCode:     core_errors.CodeAlreadyAtDestination,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, err := NavigateWarp(test.ship, test.newSystem)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.newSystem.X, ship.SystemX)
				assert.Equal(t, test.newSystem.Y, ship.SystemY)
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
		errCode     core_errors.ErrorCode
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
			err:         core_errors.ErrUnprocessableEntity,
			errCode:     core_errors.CodeInvalidShipState,
		},
		{
			name: "Invalid coordinates",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
			},
			system:      worldgen.System{},
			orbitIndex:  1,
			expectedErr: true,
			err:         core_errors.ErrNotFound,
			errCode:     core_errors.CodeInvalidCoordinates,
		},
		{
			name: "Already at destination",
			ship: model.Ship{
				Status:     model.ShipStatusOrbit,
				Location:   model.ShipLocationPlanet,
				LocationID: 1,
			},
			system:      worldgen.System{},
			orbitIndex:  1,
			expectedErr: true,
			err:         core_errors.ErrNotModified,
			errCode:     core_errors.CodeAlreadyAtDestination,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, err := NavigatePlanet(test.ship, test.system, test.orbitIndex)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, model.ShipLocationPlanet, ship.Location)
				assert.Equal(t, test.orbitIndex, ship.LocationID)
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
		errCode     core_errors.ErrorCode
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
			err:         core_errors.ErrUnprocessableEntity,
			errCode:     core_errors.CodeInvalidShipState,
		},
		{
			name: "Invalid coordinates",
			ship: model.Ship{
				Status: model.ShipStatusOrbit,
			},
			system:      worldgen.System{},
			waypointID:  0,
			expectedErr: true,
			err:         core_errors.ErrNotFound,
			errCode:     core_errors.CodeInvalidCoordinates,
		},
		{
			name: "Already at destination",
			ship: model.Ship{
				Status:     model.ShipStatusOrbit,
				Location:   model.ShipLocationWaypoint,
				LocationID: 0,
			},
			system:      worldgen.System{},
			waypointID:  0,
			expectedErr: true,
			err:         core_errors.ErrNotModified,
			errCode:     core_errors.CodeAlreadyAtDestination,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ship, cd, err := NavigateWaypoint(test.ship, test.system, test.waypointID)
			if test.expectedErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				var withCode core_errors.WithCode
				if assert.ErrorAs(t, err, &withCode) {
					assert.Equal(t, test.errCode, withCode.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, model.ShipLocationWaypoint, ship.Location)
				assert.Equal(t, test.waypointID, ship.LocationID)
				assert.True(t, cd > 0)
			}
		})
	}
}
