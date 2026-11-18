package http_server

import (
	"fmt"
	"net/http"
)

type APIVersion string

const (
	APIVersionV1 APIVersion = "v1"
	APIVersionV2 APIVersion = "v2"
)

type APIVersionRouter struct {
	*http.ServeMux

	apiVersion APIVersion
}

func NewAPIVersionRouter(apiVersion APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) AddRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.WithMiddleware())
	}
}
