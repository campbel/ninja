package ninja

import (
	"net/http"
)

// Handle can be used as an http.Handle with support for muxing HTTP Methods and Middleware
type Handle struct {
	mux *http.ServeMux
}

// New creates a new ninja handle
func New() *Handle {
	return &Handle{
		mux: http.NewServeMux(),
	}
}

func (h *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.mux.ServeHTTP(w, r)
}

// Route registers a handler and middleware at a given route
func (h *Handle) Route(pattern string, route Route, middlewares ...Middleware) *Handle {

	h.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		for _, middleware := range middlewares {
			if middleware(w, r) != nil {
				return
			}
		}

		handleFn, exists := route[r.Method]

		if exists {
			handleFn(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})

	return h
}

// Route is a map from HTTP Method to HandleFunc
type Route map[string]func(http.ResponseWriter, *http.Request)

// Middleware is the function signature for a middleware method
type Middleware func(http.ResponseWriter, *http.Request) error
