package handler

import (
	"astragalaxy/internal/schemas"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	body := &schemas.CreateUserSchema{TelegramID: 987654321, Username: "tester2"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("secret-token", sudoToken)

	res, err := app.Test(request, -1)
	assert.NoError(t, err)

	if assert.Equal(t, res.StatusCode, http.StatusCreated) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response *schemas.UserSchema

		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, int64(987654321), response.TelegramID)
		assert.Equal(t, "tester2", response.Username)

		spaceship, err := stateObj.S.FindOneSpaceship(response.Spaceships[0].ID)
		assert.NoError(t, err)

		assert.Equal(t, "initial", spaceship.Name)
		assert.Equal(t, response.ID, spaceship.UserID)
	}
}

func TestLogin(t *testing.T) {
	body := fmt.Sprintf(`{"token":"%s"}`, fmt.Sprintf("%d:%s", usr.TelegramID, userToken))

	request := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	res, err := app.Test(request, -1)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var authBody *schemas.AuthBody
		err = json.Unmarshal(body, &authBody)
		assert.NoError(t, err)

		assert.Equal(t, authBody.TokenType, "Bearer")
	}
}

func TestGetUserTokenSudo(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/token/sudo?telegram_id=%d", usr.TelegramID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("secret-token", sudoToken)

	res, err := app.Test(request, -1)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var response *schemas.UserTokenResponseSchema
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, userToken, response.Token)
	}
}

func TestGetMe(t *testing.T) {
	url := "/auth/me"
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userJwtToken))
	res, err := app.Test(request, -1)

	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, res.StatusCode) {
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var me *schemas.UserSchema
		err = json.Unmarshal(body, &me)
		assert.NoError(t, err)

		assert.NotEmpty(t, me)
		//assert.Equal(t, me.TelegramID, usr.TelegramID)
		assert.Equal(t, me, usr)
	}
}
