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

	r := mux.NewRouter()
	rest.NewRouter(r, &routes.Auth{}, &routes.User{})
	r.Use(mux.CORSMethodMiddleware(r))
	http.ListenAndServe(cfg.HTTPServer.Address, r)

	log.Info("starting restapi server")
	log.Debug("debug messages are enabled")
}
