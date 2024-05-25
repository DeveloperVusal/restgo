package app

import (
	"log/slog"
	"net/http"

	"apibgo/internal/app/instance"
	"apibgo/internal/config"
	"apibgo/internal/storage"
	"apibgo/internal/transport/rest"
	"apibgo/internal/transport/rest/middleware"
	"apibgo/internal/transport/rest/routes"

	"github.com/gorilla/mux"
)

type Instance struct {
	Config  *config.Config
	Storage *storage.Config
	Log     *slog.Logger
}

func Run() {
	// Load .env files for global
	// univenv.Load()

	instance := instance.GetInstance()

	_routes := []rest.Handler{
		&routes.Auth{Config: instance.Config, Storage: instance.Storage},
		&routes.User{
			Config:  instance.Config,
			Storage: instance.Storage,
			Middlewares: []mux.MiddlewareFunc{
				middleware.LoggingMiddleware,
			},
		},
	}

	r := mux.NewRouter()
	rest.NewRouter(r, _routes...)
	r.Use(mux.CORSMethodMiddleware(r))

	if err := http.ListenAndServe(instance.Config.HTTPServer.Address, r); err == nil {
		instance.Log.Info("starting restapi server")
	}
}
