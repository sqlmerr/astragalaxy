package v1

import (
	"astragalaxy/internal/schema"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	api := createAPI(t)
	body := &schema.CreateUser{Password: "987654321", Username: "tester2"}

	res := api.Post("/v1/auth/register", fmt.Sprintf("secret-token: %s", testSudoToken), body)

	if assert.Equal(t, http.StatusCreated, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response *schema.User

		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, "tester2", response.Username)
	}
}

func TestLoginByToken(t *testing.T) {
	api := createAPI(t)
	body := schema.AuthPayloadToken{Token: testUserToken}

	res := api.Post("/v1/auth/login/token", body)

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var authBody *schema.AuthBody
		err = json.Unmarshal(body, &authBody)
		assert.NoError(t, err)

		assert.Equal(t, authBody.TokenType, "Bearer")
	}
}

func TestLogin(t *testing.T) {
	api := createAPI(t)
	body := schema.AuthPayload{Username: testUser.Username, Password: "testPassword"}

	res := api.Post("/v1/auth/login", body)

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var authBody *schema.AuthBody
		err = json.Unmarshal(body, &authBody)
		assert.NoError(t, err)

		assert.Equal(t, authBody.TokenType, "Bearer")
	}
}

func TestGetUserTokenSudo(t *testing.T) {
	api := createAPI(t)
	url := fmt.Sprintf("/v1/auth/token/sudo?id=%s", testUser.ID.String())

	res := api.Get(url, fmt.Sprintf("secret-token: %s", testSudoToken))
	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response *schema.UserTokenResponse
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, testUserToken, response.Token)
	}
}

func TestGetMe(t *testing.T) {
	api := createAPI(t)
	url := "/v1/auth/me"
	res := api.Get(url, fmt.Sprintf("Authorization: %s", fmt.Sprintf("Bearer %s", testUserJwtToken)))

	if assert.Equal(t, http.StatusOK, res.Code) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var me *schema.User
		err = json.Unmarshal(body, &me)
		assert.NoError(t, err)

		assert.NotEmpty(t, me)
		//assert.Equal(t, me.TelegramID, testUser.TelegramID)
		assert.Equal(t, me, testUser)
	}
}
