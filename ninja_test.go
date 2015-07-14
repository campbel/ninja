package ninja_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/campbel/ninja"
)

type mockResponseWriter struct{}

func (mock mockResponseWriter) Header() http.Header             { return http.Header{} }
func (mock mockResponseWriter) Write(bytes []byte) (int, error) { return len(bytes), nil }
func (mock mockResponseWriter) WriteHeader(status int)          {}

func TestSimpleNinja(t *testing.T) {

	api := ninja.New()

	calledRoute := false
	calledMiddleware := false
	api.Route("/foo/bar",
		ninja.Route{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				calledRoute = true
			},
		},
		func(w http.ResponseWriter, r *http.Request) error {
			calledMiddleware = true
			return nil
		},
	)

	request, _ := http.NewRequest("GET", "/foo/bar", nil)

	api.ServeHTTP(mockResponseWriter{}, request)

	if !calledRoute {
		t.Error("Should have called the route method")
	}

	if !calledMiddleware {
		t.Error("Should have called the middleware method")
	}
}

func TestMiddleware(t *testing.T) {
	api := ninja.New()

	calledRoute := false
	calledMiddleware := false
	api.Route("/foo/bar",
		ninja.Route{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				calledRoute = true
			},
		},
		func(w http.ResponseWriter, r *http.Request) error {
			calledMiddleware = true
			return errors.New("")
		},
	)

	request, _ := http.NewRequest("GET", "/foo/bar", nil)

	api.ServeHTTP(mockResponseWriter{}, request)

	if calledRoute {
		t.Error("Should NOT have called the route method")
	}

	if !calledMiddleware {
		t.Error("Should have called the middleware method")
	}
}

func TestSubRoute(t *testing.T) {
	api := ninja.New()

	calledRoute1 := false
	calledRoute2 := false
	api.Route("/foo/bar",
		ninja.Route{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				calledRoute1 = true
			},
		},
	)
	api.Route("/foo/bar/baz",
		ninja.Route{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				calledRoute2 = true
			},
		},
	)
	request, _ := http.NewRequest("GET", "/foo/bar/baz", nil)

	api.ServeHTTP(mockResponseWriter{}, request)

	if calledRoute1 {
		t.Error("Should NOT have called the less specific route")
	}

	if !calledRoute2 {
		t.Error("Should have called the more specific route")
	}
}

func TestMethodMuxing(t *testing.T) {
	api := ninja.New()

	calledRoute1 := false
	calledRoute2 := false
	api.Route("/foo/bar",
		ninja.Route{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				calledRoute1 = true
			},
			"POST": func(w http.ResponseWriter, r *http.Request) {
				calledRoute2 = true
			},
		},
	)
	request, _ := http.NewRequest("POST", "/foo/bar", nil)

	api.ServeHTTP(mockResponseWriter{}, request)

	if calledRoute1 {
		t.Error("Should NOT have called the GET method")
	}

	if !calledRoute2 {
		t.Error("Should have called the POST method")
	}
}

func TestRouteCoveringSubRoutes(t *testing.T) {
	api := ninja.New()

	calledRoute := false
	api.Route("/foo/bar/",
		ninja.Route{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				calledRoute = true
			},
		},
	)
	request, _ := http.NewRequest("GET", "/foo/bar/baz", nil)

	api.ServeHTTP(mockResponseWriter{}, request)

	if !calledRoute {
		t.Error("Should have called the route method")
	}
}

func TestRouteNotCoveringSubRoutes(t *testing.T) {
	api := ninja.New()

	calledRoute := false
	api.Route("/foo/bar",
		ninja.Route{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				calledRoute = true
			},
		},
	)
	request, _ := http.NewRequest("GET", "/foo/bar/baz", nil)

	api.ServeHTTP(mockResponseWriter{}, request)

	if calledRoute {
		t.Error("Should NOT have called the route method")
	}
}
