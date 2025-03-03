package test

import (
	"bytes"
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
		var bodyReader io.Reader
		if test.Body != nil {
			bodyReader = bytes.NewReader(test.Body)
		}
		req := httptest.NewRequest(test.Method, test.Route, bodyReader)
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}

		res, err := e.app.Test(req, -1)

		assert.NoError(t, err, test.Description)
		//b, err := io.ReadAll(res.Body)
		//var arbuz interface{}
		//json.Unmarshal(b, &arbuz)
		//fmt.Println(res.StatusCode, arbuz)
		assert.Equalf(t, test.ExpectedCode, res.StatusCode, test.Description)
		if test.ExpectedError || test.BodyValidator == nil {
			continue
		}
		body, err := io.ReadAll(res.Body)
		assert.NoErrorf(t, err, test.Description)
		if test.BodyValidator != nil {
			test.BodyValidator(body)
		}
	}
}
