package handler

import (
	"astragalaxy/internal/schema"
	"astragalaxy/pkg/test"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetMyItems(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "Get items",
			Route:         "/inventory/items",
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.DataGenericResponse[[]schema.Item]
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.Len(t, b.Data, 1) {
					assert.Equal(t, "test", b.Data[0].Code)
					assert.Equal(t, testUser.ID, b.Data[0].UserID)
				}
			},
		},
	}

	testExecutor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken)})
}

func TestGetMyItemsByCode(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "Get items by code",
			Route:         "/inventory/items/test",
			Method:        http.MethodGet,
			ExpectedError: false,
			ExpectedCode:  200,
			BodyValidator: func(body []byte) {
				var b schema.DataGenericResponse[[]schema.Item]
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.Len(t, b.Data, 1) {
					assert.Equal(t, "test", b.Data[0].Code)
					assert.Equal(t, testUser.ID, b.Data[0].UserID)
				}
			},
		},
		{
			Description:   "Get items by code not found",
			Route:         "/inventory/items/notfound",
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

	testExecutor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken)})
}

func TestGetItemData(t *testing.T) {
	tests := []test.HTTPTest{
		{
			Description:   "Get item data success",
			Route:         fmt.Sprintf("/inventory/items/%s/data", testItem.ID.String()),
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
			Route:         "/inventory/items/123/data",
			Method:        http.MethodGet,
			ExpectedError: true,
			ExpectedCode:  400,
		},
		{
			Description:   "Item not found",
			Route:         "/inventory/items/34afa4f5-c0e9-49ca-8e13-7dcb731b1541/data",
			Method:        http.MethodGet,
			ExpectedError: true,
			ExpectedCode:  404,
		},
	}

	testExecutor.TestHTTP(t, tests, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken)})
}
