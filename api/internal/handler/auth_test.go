package handler

import (
	"astragalaxy/internal/schemas"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
