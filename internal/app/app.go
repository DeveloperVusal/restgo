package app

import (
	"net/http"

	"apibgo/internal/config"
	"apibgo/internal/storage"
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
	dbcfgs := storage.MustLoad()

	_routes := []rest.Handler{
		&routes.Auth{Config: cfg, Storage: dbcfgs}, &routes.User{},
	}

	r := mux.NewRouter()
	rest.NewRouter(r, _routes...)
	r.Use(mux.CORSMethodMiddleware(r))

	if err := http.ListenAndServe(cfg.HTTPServer.Address, r); err == nil {
		log.Info("starting restapi server")
	}
}
