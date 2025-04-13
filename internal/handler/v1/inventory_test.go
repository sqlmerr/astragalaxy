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

func TestGetMyItems(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	tests := []test.HTTPTest{
		{
			Description:   "Get items",
			Route:         "/v1/inventory/items",
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.DataGenericResponse[[]schema.Item]
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.Len(t, b.Data, 1) {
					assert.Equal(t, "test", b.Data[0].Code)
					assert.Equal(t, testAstral.ID, b.Data[0].AstralID)
				}
			},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestGetMyItemsByCode(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	tests := []test.HTTPTest{
		{
			Description:   "Get items by code",
			Route:         "/v1/inventory/items/test",
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.DataGenericResponse[[]schema.Item]
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.Len(t, b.Data, 1) {
					assert.Equal(t, "test", b.Data[0].Code)
					assert.Equal(t, testAstral.ID, b.Data[0].AstralID)
				}
			},
		},
		{
			Description:   "Get items by code not found",
			Route:         "/v1/inventory/items/notfound",
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.DataResponse
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.Len(t, b.Data, 0)
			},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestGetItemData(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)
	anotherUuid := uuid.New().String()

	fmt.Println("asdasdasjkklajsdf, ", anotherUuid == testItem.ID.String(), anotherUuid, testItem.ID.String())
	tests := []test.HTTPTest{
		{
			Description:   "Get item data success",
			Route:         fmt.Sprintf("/v1/inventory/items/%s/data", testItem.ID.String()),
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
			Route:         "/v1/inventory/items/123/data",
			Method:        http.MethodGet,
			ExpectedError: true,
			ExpectedCode:  400,
		},
		{
			Description:   "Item not found",
			Route:         fmt.Sprintf("/v1/inventory/items/%s/data", anotherUuid),
			Method:        http.MethodGet,
			ExpectedError: true,
			ExpectedCode:  404,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}
