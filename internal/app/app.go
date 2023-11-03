package app

import (
	"net/http"

	"apibgo/internal/config"
	"apibgo/internal/transport/rest"
	"apibgo/internal/transport/rest/routes"

	"apibgo/pkg/logger"
	"apibgo/pkg/univenv"

	"github.com/gorilla/mux"
)

func Run() {
	// Load .env files
	univenv.Load()

	cfg := config.MustLoad()
	log := logger.Setup(cfg.Env)

	// db := storage.MustLoad()

	log.Info("starting restapi server")

	_routes := []rest.Handler{
		&routes.Auth{}, &routes.User{},
	}

	r := mux.NewRouter()
	rest.NewRouter(r, _routes...)
	r.Use(mux.CORSMethodMiddleware(r))
	http.ListenAndServe(cfg.HTTPServer.Address, r)
}
