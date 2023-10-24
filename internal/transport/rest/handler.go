package rest

import (
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
