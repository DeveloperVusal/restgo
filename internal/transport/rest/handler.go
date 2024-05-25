package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	NewHandler(*mux.Router)
}

func NewRouter(r *mux.Router, routes ...Handler) {

	for _, route := range routes {
		route.NewHandler(r)
	}

}

func Adapt(handler http.Handler, adapters ...mux.MiddlewareFunc) http.Handler {
	// The loop is reversed so the adapters/middleware gets executed in the same
	// order as provided in the array.
	for i := len(adapters); i > 0; i-- {
		handler = adapters[i-1](handler)
	}
	return handler
}
