package v1

import (
	"astragalaxy/internal/schema"
	"astragalaxy/pkg/test"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestGetAstralBundle(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	bundle, err := testStateObj.S.GetBundleByInventory(testAstralInventory.ID)
	assert.NoError(t, err)

	tests := []test.HTTPTest{
		{
			Description:   "success",
			Route:         "/v1/inventory/my/bundle",
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(b []byte) {
				var body schema.Bundle
				err = json.Unmarshal(b, &body)
				assert.NoError(t, err)
				assert.Equal(t, body.ID, bundle.ID)
				assert.Equal(t, body.InventoryID, bundle.InventoryID)
			},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestGetBundle(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	bundleAstral, err := testStateObj.S.GetBundleByInventory(testAstralInventory.ID)
	assert.NoError(t, err)

	bundleSpaceship, err := testStateObj.S.GetBundleByInventory(testSpaceshipInventory.ID)
	assert.NoError(t, err)

	tests := []test.HTTPTest{
		{
			Description:   "Successfully got astral bundle",
			Route:         fmt.Sprintf("/v1/inventory/%s/bundle", testAstralInventory.ID.String()),
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(b []byte) {
				var body schema.Bundle
				err = json.Unmarshal(b, &body)
				assert.NoError(t, err)
				assert.Equal(t, body.ID, bundleAstral.ID)
				assert.Equal(t, body.InventoryID, bundleAstral.InventoryID)
			},
		},
		{
			Description:   "Successfully got spaceship bundle",
			Route:         fmt.Sprintf("/v1/inventory/%s/bundle", testSpaceshipInventory.ID.String()),
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(b []byte) {
				var body schema.Bundle
				err = json.Unmarshal(b, &body)
				assert.NoError(t, err)
				assert.Equal(t, body.ID, bundleSpaceship.ID)
				assert.Equal(t, body.InventoryID, bundleSpaceship.InventoryID)
			},
		},
		{
			Description:   "Not Found",
			Route:         fmt.Sprintf("/v1/inventory/%s/bundle", uuid.New().String()),
			Method:        http.MethodGet,
			ExpectedError: true,
			ExpectedCode:  404,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}
