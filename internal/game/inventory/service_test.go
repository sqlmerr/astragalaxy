package inventory_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sqlmerr/astragalaxy/internal/data"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/data/redis"
	"github.com/sqlmerr/astragalaxy/internal/data/registry"
	errs "github.com/sqlmerr/astragalaxy/internal/errors"
	"github.com/sqlmerr/astragalaxy/internal/game"
	"github.com/sqlmerr/astragalaxy/internal/game/inventory"
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

	code := m.Run()
	pgCleanup()
	redisCleanup()
	os.Exit(code)
}

func TestTrash(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	service := inventory.NewService(game.Config{}, testStore, *testGameData)

	t.Run("success", func(t *testing.T) {
		user := testutil.CreateUser(t, testStore)
		agent, inv := testutil.CreateAgent(t, testStore, user.ID)
		item := testutil.AddItem(t, testStore, inv.ID, model.ItemPortablePrinter)
		res := testutil.AddResource(t, testStore, inv.ID, "iron", 100)

		cooldown, err := service.Trash(
			t.Context(),
			agent.ID,
			inv.ID,
			[]uuid.UUID{item.ID},
			map[model.ResourceType]int{res.ResourceType: 50},
		)

		require.NoError(t, err)
		assert.Equal(t, time.Second*5, cooldown.Duration)
		assert.Equal(t, "trash", cooldown.Action)

		newInvItemsCount, err := testStore.Inventories().GetInventoryItemCount(t.Context(), inv.ID)
		require.NoError(t, err)
		assert.Equal(t, 0, newInvItemsCount)

		newResource, err := testStore.Inventories().GetResource(t.Context(), inv.ID, res.ResourceType)
		require.NoError(t, err)
		assert.Equal(t, 50, newResource.Amount)

		t.Run("cooldown", func(t *testing.T) {
			_, err := service.Trash(
				t.Context(),
				agent.ID,
				inv.ID,
				nil,
				nil,
			)
			assert.Error(t, err)
			testutil.AssertErrorCode(t, err, errs.CodeAgentInCooldown)
		})
	})

	t.Run("can't access inventory", func(t *testing.T) {
		user := testutil.CreateUser(t, testStore)
		agent, _ := testutil.CreateAgent(t, testStore, user.ID)
		_, inv2 := testutil.CreateAgent(t, testStore, user.ID)

		_, err := service.Trash(t.Context(), agent.ID, inv2.ID, nil, nil)
		assert.Error(t, err)
		testutil.AssertErrorCode(t, err, errs.CodeAccessDenied)
	})

	t.Run("not enough resources", func(t *testing.T) {
		user := testutil.CreateUser(t, testStore)
		agent, inv := testutil.CreateAgent(t, testStore, user.ID)
		res := testutil.AddResource(t, testStore, inv.ID, "crystal", 1000)

		_, err := service.Trash(
			t.Context(),
			agent.ID,
			inv.ID,
			nil,
			map[model.ResourceType]int{res.ResourceType: res.Amount + 1},
		)
		assert.Error(t, err)
		testutil.AssertErrorCode(t, err, errs.CodeNotEnoughResources)
	})

	t.Run("item not found", func(t *testing.T) {
		user := testutil.CreateUser(t, testStore)
		agent, inv := testutil.CreateAgent(t, testStore, user.ID)

		_, err := service.Trash(
			t.Context(),
			agent.ID,
			inv.ID,
			[]uuid.UUID{uuid.New()},
			nil,
		)
		assert.Error(t, err)
		testutil.AssertErrorCode(t, err, errs.CodeItemNotFound)
	})
}
