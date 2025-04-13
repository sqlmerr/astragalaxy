package test

import (
	"bytes"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

type Executor struct {
	app humatest.TestAPI
}

func New(app humatest.TestAPI) *Executor {
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

		res := e.app.Do(test.Method, test.Route)

		assert.Equalf(t, test.ExpectedCode, res.Code, test.Description)
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
