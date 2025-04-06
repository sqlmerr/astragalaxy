package handler

import (
	"astragalaxy/internal/schema"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	body := &schema.CreateUser{Password: "987654321", Username: "tester2"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("secret-token", testSudoToken)

	res, err := testApp.Test(request, -1)
	assert.NoError(t, err)

	if assert.Equal(t, res.StatusCode, http.StatusCreated) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response *schema.User

		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, "tester2", response.Username)

		spaceship, err := testStateObj.S.FindOneSpaceship(response.Spaceships[0].ID)
		assert.NoError(t, err)

		assert.Equal(t, "initial", spaceship.Name)
		assert.Equal(t, response.ID, spaceship.UserID)
	}
}

func TestLoginByToken(t *testing.T) {
	body := fmt.Sprintf(`{"token":"%s"}`, testUserToken)

	request := httptest.NewRequest(http.MethodPost, "/auth/login/token", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	res, err := testApp.Test(request, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var authBody *schema.AuthBody
		err = json.Unmarshal(body, &authBody)
		assert.NoError(t, err)

		assert.Equal(t, authBody.TokenType, "Bearer")
	}
}

func TestLogin(t *testing.T) {
	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, testUser.Username, "testPassword")

	request := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	res, err := testApp.Test(request, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var authBody *schema.AuthBody
		err = json.Unmarshal(body, &authBody)
		assert.NoError(t, err)

		assert.Equal(t, authBody.TokenType, "Bearer")
	}
}

func TestGetUserTokenSudo(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/token/sudo?id=%s", testUser.ID.String()), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("secret-token", testSudoToken)

	res, err := testApp.Test(request, -1)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response *schema.UserTokenResponse
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, testUserToken, response.Token)
	}
}

func TestGetMe(t *testing.T) {
	url := "/auth/me"
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testUserJwtToken))
	res, err := testApp.Test(request, -1)

	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
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
