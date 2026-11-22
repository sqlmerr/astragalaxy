package mail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sqlmerr/astragalaxy/internal/data"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/data/redis"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/facilities"
	"github.com/sqlmerr/astragalaxy/internal/game/mail"
	"github.com/sqlmerr/astragalaxy/internal/game/worldgen"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"github.com/sqlmerr/astragalaxy/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testPool     postgres_pool.Pool
	testRedis    redis.Redis
	testStore    data.Store
	testGameData *registry.GameData
	testWorldGen *worldgen.WorldGen
)

func TestMain(m *testing.M) {
	if err := os.Chdir(testutil.ProjectRoot()); err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, pgCleanup, err := testutil.SetupPostgres(ctx)
	if err != nil {
		panic(err)
	}
	testPool = pool

	rdb, redisCleanup, err := testutil.SetupRedis(ctx)
	if err != nil {
		pgCleanup()
		panic(err)
	}
	testRedis = rdb

	testStore = testutil.SetupStore(pool, rdb)
	gameData, err := registry.LoadGameData(registry.LoadConfigMust())
	if err != nil {
		pgCleanup()
		redisCleanup()
		panic(err)
	}
	testGameData = gameData

	testWorldGen = worldgen.New(testGameData, 0)

	code := m.Run()
	pgCleanup()
	redisCleanup()
	os.Exit(code)
}

func TestClaimPackage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	facilityService := facilities.New(testStore, *testGameData, *testWorldGen)
	service := mail.New(testStore, game.LoadConfigMust(), facilityService)
	user := testutil.CreateUser(t, testStore)
	system, err := testWorldGen.FindSpawnSystem()
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		agent, _ := testutil.CreateAgent(t, testStore, user.ID)
		testutil.CreateShip(t, testStore, agent.ID, true, model.ShipStatusDocked, model.ShipCoords{
			Location:   model.ShipLocationWaypoint,
			LocationID: 0,
			SystemX:    system.X,
			SystemY:    system.Y,
		})
		pkg, _ := testutil.CreateMailPackage(
			t,
			testStore,
			nil,
			agent.ID,
			[]model.ItemType{model.ItemContract},
			map[model.ResourceType]int{"iron": 5},
			time.Second*-1,
		)

		cooldown, err := service.ClaimPackage(t.Context(), agent, pkg.ID)
		require.NoError(t, err)

		assert.Equal(t, 10*time.Second, cooldown.Duration)
		assert.Equal(t, "claim_mail_package", cooldown.Action)

		newInvItems, err := testStore.Inventories().GetInventoryItems(t.Context(), agent.InventoryID)
		require.NoError(t, err)
		require.Len(t, newInvItems, 1)
		assert.Equal(t, model.ItemContract, newInvItems[0].ItemType)

		newInvResources, err := testStore.Inventories().GetInventoryResources(t.Context(), agent.InventoryID)
		require.NoError(t, err)
		require.Len(t, newInvResources, 1)
		assert.Equal(t, 5, newInvResources[0].Amount)
		assert.Equal(t, model.ResourceType("iron"), newInvResources[0].ResourceType)
	})

	t.Run("inventory is full", func(t *testing.T) {
		agent, inv := testutil.CreateAgent(t, testStore, user.ID)
		testutil.CreateShip(t, testStore, agent.ID, true, model.ShipStatusDocked, model.ShipCoords{
			Location:   model.ShipLocationWaypoint,
			LocationID: 0,
			SystemX:    system.X,
			SystemY:    system.Y,
		})
		pkg, _ := testutil.CreateMailPackage(
			t,
			testStore,
			nil,
			agent.ID,
			[]model.ItemType{model.ItemContract},
			map[model.ResourceType]int{"iron": inv.MaxResourceVolume + 5},
			time.Second*-1,
		)

		_, err := service.ClaimPackage(t.Context(), agent, pkg.ID)
		require.Error(t, err)
		testutil.AssertErrorCode(t, err, errs.CodeInventoryIsFull)
	})

	t.Run("package is not delivered", func(t *testing.T) {
		agent, _ := testutil.CreateAgent(t, testStore, user.ID)
		testutil.CreateShip(t, testStore, agent.ID, true, model.ShipStatusDocked, model.ShipCoords{
			Location:   model.ShipLocationWaypoint,
			LocationID: 0,
			SystemX:    system.X,
			SystemY:    system.Y,
		})
		pkg, _ := testutil.CreateMailPackage(
			t,
			testStore,
			nil,
			agent.ID,
			[]model.ItemType{model.ItemContract},
			map[model.ResourceType]int{"iron": 5},
			time.Second*5,
		)

		_, err := service.ClaimPackage(t.Context(), agent, pkg.ID)
		require.Error(t, err)
		testutil.AssertErrorCode(t, err, errs.CodeCannotClaimPackage)
	})

	t.Run("facility not found", func(t *testing.T) {
		agent, _ := testutil.CreateAgent(t, testStore, user.ID)
		testutil.CreateShip(t, testStore, agent.ID, true, model.ShipStatusOrbit, model.ShipCoords{
			Location:   model.ShipLocationNone,
			LocationID: 0,
			SystemX:    system.X,
			SystemY:    system.Y,
		})
		pkg, _ := testutil.CreateMailPackage(
			t,
			testStore,
			nil,
			agent.ID,
			[]model.ItemType{model.ItemContract},
			map[model.ResourceType]int{"iron": 5},
			time.Second*-1,
		)

		_, err := service.ClaimPackage(t.Context(), agent, pkg.ID)
		require.Error(t, err)
		testutil.AssertErrorCode(t, err, errs.CodeFacilityNotFound)
	})
}
