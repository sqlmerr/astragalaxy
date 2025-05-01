package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/pkg/test"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetMyAstralInventory(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	tests := []test.HTTPTest{
		{
			Description:   "Get astral inventory",
			Route:         "/v1/inventory/items/my",
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.Inventory
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.Equal(t, testAstral.ID, b.HolderID) && assert.Len(t, b.Items, 1) && assert.Equal(t, b.ID, testAstralInventory.ID) {
					assert.Equal(t, "test", b.Items[0].Code)
					assert.Equal(t, b.ID, b.Items[0].InventoryID)
				}
			},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestGetHolderInventory(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	tests := []test.HTTPTest{
		{
			Description:   "Get astral holder inventory",
			Route:         fmt.Sprintf("/v1/inventory/items/astral/%s", testAstral.ID.String()),
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.Inventory
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.Equal(t, testAstral.ID, b.HolderID) && assert.Len(t, b.Items, 1) && assert.Equal(t, b.ID, testAstralInventory.ID) {
					assert.Equal(t, "test", b.Items[0].Code)
					assert.Equal(t, b.ID, b.Items[0].InventoryID)
				}
			},
		},
		{
			Description:   "Get spaceship holder inventory",
			Route:         fmt.Sprintf("/v1/inventory/items/spaceship/%s", testSpaceship.ID.String()),
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.Inventory
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				assert.Equal(t, testSpaceship.ID, b.HolderID)
				assert.Len(t, b.Items, 0)
				assert.Equal(t, b.ID, testSpaceshipInventory.ID)
			},
		},
		{
			Description:   "Inventory not found",
			Route:         fmt.Sprintf("/v1/inventory/items/notfound/%s", uuid.New().String()),
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  404,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestGetItemData(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	anotherUuid := uuid.New().String()
	fmt.Println("arar", testItem.ID)

	tests := []test.HTTPTest{
		{
			Description:   "Get item data success",
			Route:         fmt.Sprintf("/v1/inventory/items/my/%s/data", testItem.ID.String()),
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.ItemDataResponse
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.Equal(t, "123", b.Data["test"])
			},
		},
		{
			Description:   "Invalid uuid",
			Route:         "/v1/inventory/items/my/123/data",
			Method:        http.MethodGet,
			ExpectedError: true,
			ExpectedCode:  400,
		},
		{
			Description:   "Item not found",
			Route:         fmt.Sprintf("/v1/inventory/items/my/%s/data", anotherUuid),
			Method:        http.MethodGet,
			ExpectedError: true,
			ExpectedCode:  404,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}
