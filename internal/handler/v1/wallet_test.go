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

func cleanTestWallets(t *testing.T) {
	wallets, err := testStateObj.S.GetAstralWallets(testAstral.ID)
	assert.NoError(t, err)
	assert.Len(t, wallets, 2)

	for _, w := range wallets {
		if w.ID != testWallet.ID {
			err = testStateObj.S.DeleteWallet(w.ID)
			assert.NoError(t, err)
		}
	}
}

func TestGetMyWallets(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	tests := []test.HTTPTest{
		{
			Description:   "Should return 200 OK",
			Route:         "/v1/wallets/my",
			ExpectedCode:  200,
			ExpectedError: false,
			Method:        http.MethodGet,
			BodyValidator: func(body []byte) {
				var b schema.DataGenericResponse[[]schema.Wallet]
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.Len(t, b.Data, 1) {
					assert.Equal(t, testWallet.ID, b.Data[0].ID)
					assert.Equal(t, testWallet.Name, b.Data[0].Name)
					assert.Equal(t, testAstral.ID, b.Data[0].AstralID)
					assert.Equal(t, testWallet.Locked, b.Data[0].Locked)
					assert.Equal(t, int64(1000), b.Data[0].Units)
					assert.Equal(t, int64(10), b.Data[0].Quarks)
				}
			},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestGetWalletByID(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	tests := []test.HTTPTest{
		{
			Description:   "Should return 200 OK",
			Route:         fmt.Sprintf("/v1/wallets/%s", testWallet.ID.String()),
			ExpectedCode:  200,
			ExpectedError: false,
			Method:        http.MethodGet,
			BodyValidator: func(body []byte) {
				var b schema.Wallet
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				if assert.NotEmpty(t, b) {
					assert.Equal(t, testWallet.ID, b.ID)
					assert.Equal(t, testWallet.Name, b.Name)
					assert.Equal(t, testAstral.ID, b.AstralID)
					assert.Equal(t, testWallet.Locked, b.Locked)
					assert.Equal(t, int64(1000), b.Units)
					assert.Equal(t, int64(10), b.Quarks)
				}
			},
		},
		{
			Description:   "Should return 404 Not Found",
			Route:         fmt.Sprintf("/v1/wallets/%s", uuid.New().String()),
			ExpectedCode:  404,
			ExpectedError: true,
			Method:        http.MethodGet,
		},
		{
			Description:   "Should return 400 Bad Request. Invalid uuid",
			Route:         "/v1/wallets/123123",
			ExpectedCode:  400,
			ExpectedError: true,
			Method:        http.MethodGet,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestCreateMyWallet(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	tests := []test.HTTPTest{
		{
			Description:   "Should return 201 Created",
			Route:         "/v1/wallets/my",
			ExpectedCode:  201,
			ExpectedError: false,
			Method:        http.MethodPost,
			Body:          schema.CreateWallet{Name: "testWallet2"},
			BodyValidator: func(body []byte) {
				var b schema.Wallet
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)

				assert.Equal(t, "testWallet2", b.Name)
				assert.Equal(t, testAstral.ID, b.AstralID)
			},
			AfterRequest: func() {
				cleanTestWallets(t)
			},
		},
		{
			Description:   "Should return 400 Bad Request. Too many wallets",
			Route:         "/v1/wallets/my",
			ExpectedCode:  400,
			ExpectedError: true,
			Method:        http.MethodPost,
			Body:          schema.CreateWallet{Name: "testWallet4"},
			BeforeRequest: func() {
				for _, i := range []int{2, 3} {
					_, err := testStateObj.S.CreateWallet(schema.CreateWallet{Name: fmt.Sprintf("testWallet%d", i)}, testAstral.ID)
					assert.NoError(t, err)
				}
			},
			AfterRequest: func() {
				cleanTestWallets(t)
			},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestLockMyWallet(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	tests := []test.HTTPTest{
		{
			Description:   "Should return 200 OK",
			Route:         fmt.Sprintf("/v1/wallets/my/%s/lock", testWallet.ID.String()),
			ExpectedCode:  200,
			ExpectedError: false,
			Method:        http.MethodPatch,
			BodyValidator: func(body []byte) {
				var b schema.OkResponse
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.True(t, b.Ok)
				assert.Equal(t, 1, b.CustomStatusCode)
			},
			AfterRequest: func() {
				wallet, err := testStateObj.S.GetWallet(testWallet.ID)
				assert.NoError(t, err)
				assert.True(t, wallet.Locked)
			},
		},
		{
			Description:   "Should return 404 Not Found",
			Route:         fmt.Sprintf("/v1/wallets/my/%s/lock", uuid.New().String()),
			ExpectedCode:  404,
			ExpectedError: true,
			Method:        http.MethodPatch,
		},
		{
			Description:   "Should return 400 Bad Request. Invalid uuid",
			Route:         "/v1/wallets/my/1234/lock",
			ExpectedCode:  400,
			ExpectedError: true,
			Method:        http.MethodPatch,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestUnlockMyWallet(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	tests := []test.HTTPTest{
		{
			Description:   "Should return 200 OK",
			Route:         fmt.Sprintf("/v1/wallets/my/%s/unlock", testWallet.ID.String()),
			ExpectedCode:  200,
			ExpectedError: false,
			Method:        http.MethodPatch,
			BodyValidator: func(body []byte) {
				var b schema.OkResponse
				err := json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.True(t, b.Ok)
				assert.Equal(t, 1, b.CustomStatusCode)
			},
			AfterRequest: func() {
				wallet, err := testStateObj.S.GetWallet(testWallet.ID)
				assert.NoError(t, err)
				assert.False(t, wallet.Locked)
			},
		},
		{
			Description:   "Should return 404 Not Found",
			Route:         fmt.Sprintf("/v1/wallets/my/%s/unlock", uuid.New().String()),
			ExpectedCode:  404,
			ExpectedError: true,
			Method:        http.MethodPatch,
		},
		{
			Description:   "Should return 400 Bad Request. Invalid uuid",
			Route:         "/v1/wallets/my/1234/unlock",
			ExpectedCode:  400,
			ExpectedError: true,
			Method:        http.MethodPatch,
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})
}

func TestSendCurrency(t *testing.T) {
	api := createAPI(t)
	executor := test.New(api)

	testWallet2, err := testStateObj.S.CreateWallet(schema.CreateWallet{Name: "testWallet2"}, testAstral.ID)
	assert.NoError(t, err)

	tests := []test.HTTPTest{
		{
			Description:   "Should return 200 OK",
			Route:         fmt.Sprintf("/v1/wallets/my/%s/send", testWallet.ID.String()),
			ExpectedCode:  200,
			ExpectedError: false,
			Method:        http.MethodPost,
			Body:          schema.WalletTransaction{ToWallet: testWallet2.ID, Units: 100, Quarks: 10},
			BeforeRequest: func() {
				_, err := testStateObj.S.GetWallet(testWallet.ID)
				assert.NoError(t, err)
			},
			BodyValidator: func(body []byte) {
				var b schema.OkResponse
				err = json.Unmarshal(body, &b)
				assert.NoError(t, err)
				assert.True(t, b.Ok)
				assert.Equal(t, 1, b.CustomStatusCode)

				wallet1, err := testStateObj.S.GetWallet(testWallet.ID)
				assert.NoError(t, err)

				assert.Equal(t, int64(900), wallet1.Units)
				assert.Equal(t, int64(0), wallet1.Quarks)

				wallet2, err := testStateObj.S.GetWallet(testWallet2.ID)
				assert.NoError(t, err)
				assert.Equal(t, int64(1100), wallet2.Units)
				assert.Equal(t, int64(20), wallet2.Quarks)
			},
			AfterRequest: func() {
				err = testStateObj.S.ProceedTransaction(testWallet2.ID, &schema.WalletTransaction{ToWallet: testWallet.ID, Units: 100, Quarks: 10})
				assert.NoError(t, err)
			},
		},
		{
			Description:   "Should return 400 Bad Request. Not enough money",
			Route:         fmt.Sprintf("/v1/wallets/my/%s/send", testWallet.ID.String()),
			ExpectedCode:  400,
			ExpectedError: true,
			Method:        http.MethodPost,
			Body:          schema.WalletTransaction{ToWallet: testWallet.ID, Units: 2000, Quarks: 20},
		},
		{
			Description:   "Should return 400 Bad Request. Invalid uuid",
			Route:         "/v1/wallets/my/123/send",
			ExpectedCode:  400,
			ExpectedError: true,
			Method:        http.MethodPost,
			Body:          schema.WalletTransaction{ToWallet: testWallet.ID, Units: 100, Quarks: 10},
		},
		{
			Description:   "Should return 404 Not found",
			Route:         fmt.Sprintf("/v1/wallets/my/%s/send", uuid.New().String()),
			ExpectedCode:  404,
			ExpectedError: true,
			Method:        http.MethodPost,
			Body:          schema.WalletTransaction{ToWallet: testWallet.ID, Units: 100, Quarks: 10},
		},
	}

	executor.TestHTTP(t, tests, map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", testUserJwtToken), "X-Astral-ID": testAstral.ID.String()})

	err = testStateObj.S.DeleteWallet(testWallet2.ID)
	assert.NoError(t, err)
}
