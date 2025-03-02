package test

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

type Executor struct {
	app *fiber.App
}

func New(app *fiber.App) *Executor {
	return &Executor{app: app}
}

func (e *Executor) TestHTTP(t *testing.T, tests []HTTPTest, headers ...map[string]string) {
	for _, test := range tests {
		req := httptest.NewRequest(test.Method, test.Route, nil)
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}

		res, err := e.app.Test(req, -1)

		assert.NoError(t, err, test.Description)
		assert.Equalf(t, test.ExpectedCode, res.StatusCode, test.Description)
		if test.ExpectedError || test.ExpectedBodyKeys == nil {
			continue
		}
		body, err := io.ReadAll(res.Body)
		assert.NoErrorf(t, err, test.Description)
		var bodyKeys map[string]interface{}
		err = json.Unmarshal(body, &bodyKeys)
		assert.NoErrorf(t, err, test.Description)
		for k, v := range test.ExpectedBodyKeys {
			assert.Equalf(t, v, bodyKeys[k], test.Description)
		}
	}
}
