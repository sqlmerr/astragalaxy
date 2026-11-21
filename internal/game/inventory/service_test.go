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
	user := testutil.CreateUser(t, testStore)
	agent, inv := testutil.CreateAgent(t, user, testStore)
	item := testutil.AddItem(t, inv, testStore, model.ItemPortablePrinter)

	cooldown, err := service.Trash(t.Context(), agent.ID, inv.ID, []uuid.UUID{item.ID}, nil)

	require.NoError(t, err)
	assert.Equal(t, time.Second*5, cooldown.Duration)
	assert.Equal(t, "trash", cooldown.Action)
}
