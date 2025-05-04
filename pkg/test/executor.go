package test

import (
	"fmt"
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
		var args []any
		for _, header := range headers {
			for k, v := range header {
				args = append(args, fmt.Sprintf("%s: %s", k, v))
			}
		}

		var res *httptest.ResponseRecorder
		if test.Body != nil {
			args = append(args, test.Body)
		}

		if test.BeforeRequest != nil {
			test.BeforeRequest()
		}

		if test.Body == nil {
			res = e.app.Do(test.Method, test.Route, args...)
		} else {
			res = e.app.Do(test.Method, test.Route, args...)
		}

		assert.Equalf(t, test.ExpectedCode, res.Code, test.Description)
		if test.ExpectedError || test.BodyValidator == nil {
			continue
		}
		body, err := io.ReadAll(res.Body)
		assert.NoErrorf(t, err, test.Description)
		if test.BodyValidator != nil {
			test.BodyValidator(body)
		}

		if test.AfterRequest != nil {
			test.AfterRequest()
		}
	}
}
