package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
}

func (u *User) NewHandler(r *mux.Router) {
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User"))
	}).Methods(http.MethodPost)
}
