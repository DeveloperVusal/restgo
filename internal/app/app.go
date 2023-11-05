package app

import (
	"net/http"
	"os"

	"apibgo/internal/config"
	"apibgo/internal/storage"
	"apibgo/internal/storage/pgsql"
	"apibgo/internal/transport/rest"
	"apibgo/internal/transport/rest/routes"

	"apibgo/pkg/logger"
	"apibgo/pkg/logger/feature/slog"
	"apibgo/pkg/univenv"

	"github.com/gorilla/mux"
)

func Run() {
	// Load .env files
	univenv.Load()

	cfg := config.MustLoad()
	log := logger.Setup(cfg.Env)

	dbcfgs := storage.MustLoad()
	pg, err := pgsql.New(dbcfgs, "master")

	if err != nil {
		log.Error("failed to init storage", slog.Err(err))
		os.Exit(1)
	}
	_ = pg

	log.Info("starting database")
	log.Info("starting restapi server")

	_routes := []rest.Handler{
		&routes.Auth{}, &routes.User{},
	}

	r := mux.NewRouter()
	rest.NewRouter(r, _routes...)
	r.Use(mux.CORSMethodMiddleware(r))
	http.ListenAndServe(cfg.HTTPServer.Address, r)
}
