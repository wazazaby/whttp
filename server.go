package whttp

import (
	"net/http"
)

// Middleware represents a middleware handler.
type Middleware func(http.Handler) http.Handler

// ServeMux is a wrapper around [http.ServeMux].
// It provides a way to set global middlewares as well as per handler middlewares.
type ServeMux struct {
	mux               *http.ServeMux
	globalMiddlewares []Middleware
}

// NewServeMux creates a new [ServeMux].
func NewServeMux() *ServeMux {
	return &ServeMux{mux: http.NewServeMux()}
}

// Use sets a list of global middlewares that will be applied
// to every handler registered the the [ServeMux].
func (m *ServeMux) Use(middlewares ...Middleware) {
	m.globalMiddlewares = append(m.globalMiddlewares, middlewares...)
}

// Handle registers the handler for the given pattern.
// If the given pattern conflicts, with one that is already registered, Handle panics.
func (m *ServeMux) Handle(pattern string, handler http.Handler, middlewares ...Middleware) {
	middlewares = append(m.globalMiddlewares, middlewares...)
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	m.mux.Handle(pattern, handler)
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

var (
	_ http.Handler = &ServeMux{}
)
