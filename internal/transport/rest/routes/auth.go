package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Auth struct {
}

func (a *Auth) NewHandler(r *mux.Router) {
	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})
	r.Methods(http.MethodPost)
}
